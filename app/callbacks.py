# import dash IO and graph objects
from dash.dependencies import Input, Output
# Plotly graph objects to render graph plots
import plotly.express as px
# Import dash html, bootstrap components, and tables for datatables
from dash import dcc, html, dash_table
import dash_bootstrap_components as dbc

# Import app
from app import app

# Import custom data.py
import data

circuits_df = data.circuits_df


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

# Callback to circuit bar chart, takes data request from dropdown
@app.callback(
    Output('circuits-bar', 'children'),
    Input("circuits-curves", "value"),
    Input("circuits-backends", "value"),
    Input("circuits-frameworks", "value"),
    Input("circuits-circuit", "value"),
    Input("circuits-metric", "value"),
    Input("circuits-input-dropdown", "value"),)
def update_bar_chart(
        curves_options, backends_options, framework_options, 
        circuit_option, metric_option, circuit_input
    ):
    ndf = circuits_df[
        (circuits_df['circuit'] == circuit_option) & 
        (circuits_df['curve'].isin(curves_options)) &
        (circuits_df['framework'].isin(framework_options)) &
        (circuits_df['backend'].isin(backends_options)) &
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

# Circuit line chart
@app.callback(
    Output('circuits-line', 'children'),
     Input("circuits-curves", "value"),
     Input("circuits-backends", "value"),
     Input("circuits-frameworks", "value"),
     Input("circuits-circuit", "value"),
     Input("circuits-metric", "value"),)
def update_line_chart(
        curves_options, backends_options, framework_options, 
        circuit_option, metric_option
    ):
    ndf = circuits_df[
        (circuits_df['circuit'] == circuit_option) & 
        (circuits_df['curve'].isin(curves_options)) &
        (circuits_df['framework'].isin(framework_options)) &
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
