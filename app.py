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


circuits_df = analyse_logs(LOGS, logging.CRITICAL)

app = Dash(__name__)

curves = list(set(circuits_df['curve']))
backends = list(set(circuits_df['backend']))
frameworks = list(set(circuits_df['framework']))
circuits = list(set(circuits_df['circuit']))

app.layout = html.Div([
    html.H4('ZKP Benchmarking'),
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
        options=['time', 'ram'],
        value='time',
        multi=False
    ),
    dcc.Graph(id="graph"),
])


@app.callback(
    Output("graph", "figure"), 
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


app.run_server(debug=True)
