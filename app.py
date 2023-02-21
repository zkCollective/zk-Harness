"""
Script to render the html app using dash and plotly.

Deps:
    pip install plotly dash
"""
import logging

import plotly.express as px
from dash import Dash, dcc, html, Input, Output
from processing import analyse_logs


GITHUB_REPO = "https://github.com/XXX/YYY"
HTML_HEADER = """
<h1 style="text-align:center;font-family:Georgia;font-variant:small-caps;font-size: 70px;color:#0F1419;">
ZKP Libraries Benchmarking
</h1>
<div><div style ="float:right;font-family: Courier New, monospace;">
View Source on <a href="{}">Github</a>
</div></div>
""".format(GITHUB_REPO)
LOGS = "benchmarks"


circuits_df, arithmetics_df, ec_df = analyse_logs(LOGS, logging.CRITICAL)

app = Dash(__name__)

curves = list(set(circuits_df['curve']))
backends = list(set(circuits_df['backend']))
frameworks = list(set(circuits_df['framework']))
circuits = list(set(circuits_df['circuit']))

frameworks_arithmetics = list(set(arithmetics_df['framework']))
operation_arithmetics = list(set(arithmetics_df['operation']))
field_arithmetics = list(set(arithmetics_df['field']))

frameworks_ec = list(set(ec_df['framework']))
operation_ec = list(set(ec_df['operation']))

app.layout = html.Div([
    html.H4('ZKP Benchmarking'),

    dcc.Tabs([
        dcc.Tab(label='Circuits Benchmarks', children=[
            html.Br(),
            dcc.Dropdown(
                id="curves-dropdown",
                options=curves,
                value=curves,
                multi=True
            ),
            dcc.Dropdown(
                id="backends-dropdown",
                options=backends,
                value=backends,
                multi=True
            ),
            dcc.Dropdown(
                id="frameworks-dropdown",
                options=frameworks,
                value=frameworks,
                multi=True
            ),
            dcc.Dropdown(
                id="circuit-dropdown",
                options=circuits,
                value='cubic',
                multi=False
            ),
            dcc.Dropdown(
                id="y-axis-dropdown",
                options=['time', 'ram', 'proof'],
                value='time',
                multi=False
            ),
            dcc.Graph(id="circuits-box-graph"),
        ]),
        dcc.Tab(label='Arithmetics Benchmarks', children=[
            html.Br(),
            dcc.Dropdown(
                id="frameworks-arithmetics-dropdown",
                options=frameworks_arithmetics,
                value=frameworks_arithmetics,
                multi=True
            ),
            dcc.Dropdown(
                id="opetation-arithmetics-dropdown",
                options=operation_arithmetics,
                value=operation_arithmetics,
                multi=True
            ),
            dcc.Dropdown(
                id="field-arithmetics-dropdown",
                options=field_arithmetics,
                value='base',
                multi=False
            ),
            dcc.Dropdown(
                id="y-axis-line-dropdown",
                options=['time', 'ram'],
                value='time',
                multi=False
            ),
            dcc.Graph(id="arithmetics-line-graph"),
        ]),
        dcc.Tab(label='Elliptic Curves Benchmarks', children=[
            html.Br(),
            dcc.Dropdown(
                id="frameworks-ec-dropdown",
                options=frameworks_ec,
                value=frameworks_ec,
                multi=True
            ),
            dcc.Dropdown(
                id="opetation-ec-dropdown",
                options=operation_ec,
                value=operation_ec,
                multi=True
            ),
            dcc.Dropdown(
                id="y-axis-line-ec-dropdown",
                options=['time', 'ram'],
                value='time',
                multi=False
            ),
            dcc.Graph(id="ec-line-graph"),
            
        ]),
    ]),


])


@app.callback(
    Output("circuits-box-graph", "figure"), 
    Input("curves-dropdown", "value"),
    Input("backends-dropdown", "value"),
    Input("frameworks-dropdown", "value"),
    Input("circuit-dropdown", "value"),
    Input("y-axis-dropdown", "value"))
def update_bar_chart(curves_options, backends_options, framework_options, circuit_option, y_axis):
    ndf = circuits_df[
        (circuits_df['circuit'] == circuit_option) & 
        (circuits_df['curve'].isin(curves_options)) &
        (circuits_df['framework'].isin(framework_options)) &
        (circuits_df['backend'].isin(backends_options))]

    # Create a bar chart using Plotly
    fig = px.bar(ndf, x="curve", y=y_axis, color="operation", 
                          facet_col="framework", facet_row="backend",
                          barmode="group",
                 )
    return fig

@app.callback(
    Output("arithmetics-line-graph", "figure"),
    Input("frameworks-arithmetics-dropdown", "value"),
    Input("opetation-arithmetics-dropdown", "value"),
    Input("field-arithmetics-dropdown", "value"),
    Input("y-axis-line-dropdown", "value"))
def update_arithmetics_line_chart(framework_options, operations_options, field_option, y_axis):
    ndf = arithmetics_df[
        (arithmetics_df['framework'].isin(framework_options)) &
        (arithmetics_df['operation'].isin(operations_options)) &
        (arithmetics_df['field'] == field_option)
    ]
    # TODO sort by curve fields given a lookup
    fig = px.line(ndf, x='curve', y=y_axis, color='operation',
                           facet_col="framework"
                )
    return fig


@app.callback(
    Output("ec-line-graph", "figure"),
    Input("frameworks-ec-dropdown", "value"),
    Input("opetation-ec-dropdown", "value"),
    Input("y-axis-line-ec-dropdown", "value"))
def update_ec_line_chart(framework_options, operations_options, y_axis):
    ndf = ec_df[
        (ec_df['framework'].isin(framework_options)) &
        (ec_df['operation'].isin(operations_options))
    ]
    # TODO sort by curve fields given a lookup
    fig = px.line(ndf, x='curve', y=y_axis, color='operation',
                           facet_col="framework"
                )
    return fig


app.run_server(debug=True)
