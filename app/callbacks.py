# import dash IO and graph objects
from dash.dependencies import Input, Output, State
# Plotly graph objects to render graph plots
import plotly.express as px
# Import dash html, bootstrap components, and tables for datatables
from dash import dcc, html, dash_table
import dash_bootstrap_components as dbc

# Import app
from app import app

# Import custom data.py
import data
# from index import questions, answers

circuits_df = data.circuits_df
arithmetics_df = data.arithmetics_df
ec_df = data.ec_df

#################################### CIRCUITS #################################
# This will update the circuits input dropdown 
@app.callback(
    [Output('circuits-input-dropdown', 'options'),
    Output('circuits-input-dropdown', 'value'),],
    [Input('circuits-circuit', 'value')])
def update_circuit_dropdown(selected_circuit):
    ndf = circuits_df[circuits_df['circuit'] == selected_circuit]
    circuit_inputs = list(set(ndf['input_path']))
    circuit_input = circuit_inputs[0]
    # Return the selected input (first input of circuit), and the options
    return circuit_inputs, circuit_input

# Constraints table 
@app.callback(
    [Output('circuits-data', 'children')],
    [Input('circuits-circuit', 'value'), 
     Input("circuits-backends", "value"),
     Input("circuits-frameworks", "value"),
     Input("circuits-curves", "value"),
     Input("circuits-input-dropdown", "value")])
def update_circuit_table(selected_circuit, selected_backends, selected_frameworks, selected_curves, selected_input):
    ndf = circuits_df[
        (circuits_df['circuit'] == selected_circuit) &
        (circuits_df['backend'].isin(selected_backends)) & 
        (circuits_df['framework'].isin(selected_frameworks)) & 
        (circuits_df['curve'].isin(selected_curves)) &
        (circuits_df['input_path'] == selected_input)
    ]
    # Filter unneccessary data
    circuit_data = ndf[['circuit', 'input_path', 'framework', 'backend', 'curve', 'nb_constraints']]
    circuit_data = circuit_data.drop_duplicates()

    data_note = []
    if len(ndf) == 0:
        return [html.Div(dbc.Alert('The table content is empty given the selected options.', color='warning'),)]

    data_note.append(html.Div(dash_table.DataTable(
        data= circuit_data.to_dict('records'),
        columns= [{'name': x, 'id': x} for x in circuit_data],
        style_as_list_view=True,
        editable=False,
        style_table={
            'overflowY': 'scroll',
            'width': '100%',
            'minWidth': '100%',
        },
        style_header={
                #'backgroundColor': '#f8f5f0',
                'fontWeight': 'bold'
            },
        style_cell={
                'textAlign': 'center',
                'padding': '8px',
            },
    )))
    return data_note

# Callback to circuit bar chart, takes data request from dropdown
@app.callback(
    Output('circuits-bar', 'children'),
    Input("circuits-curves", "value"),
    Input("circuits-backends", "value"),
    Input("circuits-frameworks", "value"),
    Input("circuits-circuit", "value"),
    Input("circuits-metric", "value"),
    Input("circuits-operation", "value"),
    Input("circuits-input-dropdown", "value"),)
def update_bar_chart(
        curves_options, backends_options, framework_options, 
        circuit_option, metric_option, circuit_operations, circuit_input
    ):
    ndf = circuits_df[
        (circuits_df['circuit'] == circuit_option) & 
        (circuits_df['curve'].isin(curves_options)) &
        (circuits_df['framework'].isin(framework_options)) &
        (circuits_df['backend'].isin(backends_options)) &
        (circuits_df['operation'].isin(circuit_operations)) &
        (circuits_df['input_path'] == circuit_input)]
    
    if len(ndf) == 0:
        return [html.Div(dbc.Alert('The bar chart is empty given the selected options.', color='warning'),)]
    
    # Create a bar chart using Plotly
    fig = px.bar(ndf, x="curve", y=metric_option, color="operation", 
                          facet_col="framework", facet_row="backend",
                          barmode="group", opacity=0.8, height=800
                 )
    
    return[
            dbc.Row(dbc.Col(
                dcc.Graph(
                    id='circuits-bar-graph', 
                    figure=fig,
                    config={'displayModeBar': False}), 
                xs={'size':12, 'offset':0}, 
                sm={'size':12, 'offset':0}, 
                md={'size': 12, 'offset': 0},
                lg={'size': 12, 'offset': 0}
        ))]

# Circuit line chart
@app.callback(
    Output('circuits-line', 'children'),
     Input("circuits-curves", "value"),
     Input("circuits-backends", "value"),
     Input("circuits-frameworks", "value"),
     Input("circuits-circuit", "value"),
    Input("circuits-operation", "value"),
     Input("circuits-metric", "value"),)
def update_line_chart(
        curves_options, backends_options, framework_options, 
        circuit_option, circuits_operations, metric_option
    ):
    ndf = circuits_df[
        (circuits_df['circuit'] == circuit_option) & 
        (circuits_df['curve'].isin(curves_options)) &
        (circuits_df['framework'].isin(framework_options)) &
        (circuits_df['operation'].isin(circuits_operations)) &
        (circuits_df['backend'].isin(backends_options))]
    res = []
    if len(ndf['input_path'].drop_duplicates()) <= 1:
        res.append(html.Div(dbc.Alert('The selected circuit has only one input.', color='warning'),))
        return res
    else:
        ndf['curve-operation'] = ndf['curve'] + ' ' + ndf['operation']
        fig = px.line(ndf, x="input_path", y=metric_option, color="curve-operation", 
                      facet_col="framework", facet_row="backend",
                 )
        res.append(
                dbc.Row(dbc.Col(
                    dcc.Graph(
                        id='circuits-line-graph', 
                        figure=fig,
                        config={'displayModeBar': False}), 
                    xs={'size':12, 'offset':0}, 
                    sm={'size':12, 'offset':0}, 
                    md={'size': 12, 'offset': 0},
                    lg={'size': 12, 'offset': 0}
        )))
        return res
    return fig1

################################################################################

################################# ARITHMETICS ##################################
# This will update the circuits input dropdown 
@app.callback(
    [Output('arithmetics-input-dropdown', 'options'),
    Output('arithmetics-input-dropdown', 'value'),],
    [Input('arithmetics-operation', 'value')])
def update_arithmetics_dropdown(selected_operation):
    ndf = arithmetics_df[arithmetics_df['operation'] == selected_operation]
    operation_inputs = list(set(ndf['input_path']))
    operation_input = operation_inputs[0]
    # Return the selected input (first input of operation), and the options
    return operation_inputs, operation_input

@app.callback(
    Output('arithmetics-bar', 'children'),
    Input("arithmetics-curves", "value"),
    Input("arithmetics-fields", "value"),
    Input("arithmetics-frameworks", "value"),
    Input("arithmetics-operation", "value"),
    Input("arithmetics-metric", "value"),
    Input("arithmetics-input-dropdown", "value"),)
def update_bar_chart(
        curves_options, fields_options, framework_options, 
        operation_option, metric_option, arithmetics_input
    ):
    ndf = arithmetics_df[
        (arithmetics_df['operation'] == operation_option) & 
        (arithmetics_df['curve'].isin(curves_options)) &
        (arithmetics_df['framework'].isin(framework_options)) &
        (arithmetics_df['field'].isin(fields_options)) &
        (arithmetics_df['input_path'] == arithmetics_input)]
    
    if len(ndf) == 0:
        return [html.Div(dbc.Alert('The bar chart is empty given the selected options.', color='warning'),)]
    
    # Create a bar chart using Plotly
    fig = px.bar(ndf, x="curve", y=metric_option, color="operation", 
                          facet_col="framework", facet_row="field",
                          barmode="group", opacity=0.8, height=800
                 )
    
    return[
            dbc.Row(dbc.Col(
                dcc.Graph(
                    id='arithmetics-bar-graph', 
                    figure=fig,
                    config={'displayModeBar': False}), 
                xs={'size':12, 'offset':0}, 
                sm={'size':12, 'offset':0}, 
                md={'size': 12, 'offset': 0},
                lg={'size': 12, 'offset': 0}
        ))]

################################################################################

###################################### EC ######################################
# This will update the circuits input dropdown 
@app.callback(
    [Output('ec-input-dropdown', 'options'),
    Output('ec-input-dropdown', 'value'),],
    [Input('ec-operation', 'value')])
def update_ec_dropdown(selected_operation):
    ndf = ec_df[ec_df['operation'] == selected_operation]
    operation_inputs = list(set(ndf['input_path']))
    operation_input = operation_inputs[0]
    # Return the selected input (first input of operation), and the options
    return operation_inputs, operation_input

@app.callback(
    Output('ec-bar', 'children'),
    Input("ec-curves", "value"),
    Input("ec-frameworks", "value"),
    Input("ec-operation", "value"),
    Input("ec-metric", "value"),
    Input("ec-input-dropdown", "value"),)
def update_bar_chart(
        curves_options, framework_options, 
        operation_option, metric_option, ec_input
    ):
    ndf = ec_df[
        (ec_df['operation'] == operation_option) & 
        (ec_df['curve'].isin(curves_options)) &
        (ec_df['framework'].isin(framework_options)) &
        (ec_df['input_path'] == ec_input)]
    
    if len(ndf) == 0:
        return [html.Div(dbc.Alert('The bar chart is empty given the selected options.', color='warning'),)]
    
    # Create a bar chart using Plotly
    fig = px.bar(ndf, x="curve", y=metric_option, color="operation", 
                          facet_col="framework", 
                          barmode="group", opacity=0.8, height=800
                 )
    
    return[
            dbc.Row(dbc.Col(
                dcc.Graph(
                    id='ec-bar-graph', 
                    figure=fig,
                    config={'displayModeBar': False}), 
                xs={'size':12, 'offset':0}, 
                sm={'size':12, 'offset':0}, 
                md={'size': 12, 'offset': 0},
                lg={'size': 12, 'offset': 0}
        ))]

################################################################################
