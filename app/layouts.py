# Dash components, html, and dash tables
from dash import dcc, html

# Import Bootstrap components
import dash_bootstrap_components as dbc

# Import custom data.py
import data

circuits_df = data.circuits_df
arithmetics_df = data.arithmetics_df
ec_df = data.ec_df

circuits_curves = list(set(circuits_df['curve']))
circuits_backends = list(set(circuits_df['backend']))
circuits_frameworks = list(set(circuits_df['framework']))
circuits_circuits = list(set(circuits_df['circuit']))
circuits_default_circuit = "cubic"
circuits_metrics = ["time", "ram", "proof"]
circuits_default_metric = "time"

arithmetics_curves = list(set(arithmetics_df['curve']))
arithmetics_fields = list(set(arithmetics_df['field']))
arithmetics_frameworks = list(set(arithmetics_df['framework']))
arithmetics_operations = list(set(arithmetics_df['operation']))
arithmetics_default_operation = "add"
arithmetics_metrics = ["time", "ram"]
arithmetics_default_metric = "time"

################################# CIRCUITS #####################################
circuitMenu = html.Div([
    # Curves
    dbc.Row([dbc.Col(html.H2(children='Circuit Benchmarks'))]),
    dbc.Row(
        [
            dbc.Col(html.H4(style={'text-align': 'center'}, children='Select Curve(s):'),
                xs={'size':'auto', 'offset':0}, sm={'size':'auto', 'offset':0}, md={'size':'auto', 'offset':3},
                lg={'size':'auto', 'offset':0}, xl={'size':'auto', 'offset':0}),
            dbc.Col(dcc.Dropdown(
                #style = {'text-align': 'center', 'font-size': '18px', 'width': '210px'},
                id='circuits-curves',
                options=circuits_curves,
                value=circuits_curves,
                clearable=False,
                multi=True),
                xs={'size':'auto', 'offset':0}, sm={'size':'auto', 'offset':0}, md={'size':'auto', 'offset':0},
                lg={'size':'auto', 'offset':0}, xl={'size':'auto', 'offset':0}),
        ],
    ),
    dbc.Row(dbc.Col(html.P(style={'font-size': '16px', 'opacity': '70%'},
        children='''Select which curves to display'''
    ))),
    # Frameworks
    dbc.Row(
        [
            dbc.Col(html.H4(style={'text-align': 'center'}, children='Select Framework(s):'),
                xs={'size':'auto', 'offset':0}, sm={'size':'auto', 'offset':0}, md={'size':'auto', 'offset':3},
                lg={'size':'auto', 'offset':0}, xl={'size':'auto', 'offset':0}),
            dbc.Col(dcc.Dropdown(
                #style = {'text-align': 'center', 'font-size': '18px', 'width': '210px'},
                id='circuits-frameworks',
                options=circuits_frameworks,
                value=circuits_frameworks,
                clearable=False,
                multi=True),
                xs={'size':'auto', 'offset':0}, sm={'size':'auto', 'offset':0}, md={'size':'auto', 'offset':0},
                lg={'size':'auto', 'offset':0}, xl={'size':'auto', 'offset':0}),
        ],
    ),
    dbc.Row(dbc.Col(html.P(style={'font-size': '16px', 'opacity': '70%'},
        children='''Select which frameworks to display'''
    ))),
    # Backends
    dbc.Row(
        [
            dbc.Col(html.H4(style={'text-align': 'center'}, children='Select Backend(s):'),
                xs={'size':'auto', 'offset':0}, sm={'size':'auto', 'offset':0}, md={'size':'auto', 'offset':3},
                lg={'size':'auto', 'offset':0}, xl={'size':'auto', 'offset':0}),
            dbc.Col(dcc.Dropdown(
                #style = {'text-align': 'center', 'font-size': '18px', 'width': '210px'},
                id='circuits-backends',
                options=circuits_backends,
                value=circuits_backends,
                clearable=False,
                multi=True),
                xs={'size':'auto', 'offset':0}, sm={'size':'auto', 'offset':0}, md={'size':'auto', 'offset':0},
                lg={'size':'auto', 'offset':0}, xl={'size':'auto', 'offset':0}),
        ],
    ),
    dbc.Row(dbc.Col(html.P(style={'font-size': '16px', 'opacity': '70%'},
        children='''Select which backend to display'''
    ))),
    dbc.Row(
        [
            dbc.Col(html.H4(style={'text-align': 'center'}, children='Select Circuit:'),
                xs={'size':'auto', 'offset':0}, sm={'size':'auto', 'offset':0}, md={'size':'auto', 'offset':3},
                lg={'size':'auto', 'offset':0}, xl={'size':'auto', 'offset':0}),
            dbc.Col(dcc.Dropdown(
                style = {'text-align': 'center', 'font-size': '18px', 'width': '210px'},
                id='circuits-circuit',
                options=circuits_circuits,
                value=circuits_default_circuit,
                clearable=False),
                xs={'size':'auto', 'offset':0}, sm={'size':'auto', 'offset':0}, md={'size':'auto', 'offset':0},
                lg={'size':'auto', 'offset':0}, xl={'size':'auto', 'offset':0}),

            dbc.Col(html.H4(style={'text-align': 'center', 'justify-self': 'right'}, children='Select Metric:'),
                xs={'size':'auto', 'offset':0}, sm={'size':'auto', 'offset':0}, md={'size':'auto', 'offset':3},
                lg={'size':'auto', 'offset':0}, xl={'size':'auto', 'offset':1}),
            dbc.Col(dcc.Dropdown(
                style = {'text-align': 'center', 'font-size': '18px', 'width': '210px'},
                id='circuits-metric',
                options=circuits_metrics,
                value=circuits_default_metric,
                clearable=False),
                xs={'size':'auto', 'offset':0}, sm={'size':'auto', 'offset':0}, md={'size':'auto', 'offset':0},
                lg={'size':'auto', 'offset':0}, xl={'size':'auto', 'offset':0}),
        ],
    ),
    html.Br(),
    dbc.Row(
        [
            dbc.Col(html.H4(style={'text-align': 'center'}, children='Select Input:'),
                xs={'size':'auto', 'offset':0}, sm={'size':'auto', 'offset':0}, md={'size':'auto', 'offset':3},
                lg={'size':'auto', 'offset':0}, xl={'size':'auto', 'offset':0}),
            dbc.Col(dcc.Dropdown(
                style = {'text-align': 'center', 'font-size': '18px', 'width': '500px'},
                id='circuits-input-dropdown',
                clearable=False),
                xs={'size':'auto', 'offset':0}, sm={'size':'auto', 'offset':0}, md={'size':'auto', 'offset':0},
                lg={'size':'auto', 'offset':0}, xl={'size':'auto', 'offset':0})
    ])
], className='menu')

circuitLayout = html.Div([
    # Cicrcuit Constraint Table
    dbc.Row(dbc.Col(html.Div(id='circuits-data'), xs={'size':'auto', 'offset':0}, sm={'size':'auto', 'offset':0}, md={'size':7, 'offset':0}, lg={'size':'auto', 'offset':0},
            xl={'size':10, 'offset':0}),justify="center"),
    # Bar Chart of Benchmarks
    dbc.Row(dbc.Col(html.Div(id='circuits-bar'), xs={'size':'auto', 'offset':0}, sm={'size':'auto', 'offset':0}, md={'size':7, 'offset':0}, lg={'size':'auto', 'offset':0},
            xl={'size':10, 'offset':0}),justify="center"),
    # Line Chart
    dbc.Row(dbc.Col(html.Div(id='circuits-line'), xs={'size':'auto', 'offset':0}, sm={'size':'auto', 'offset':0}, md={'size':7, 'offset':0}, lg={'size':'auto', 'offset':0},
            xl={'size':10, 'offset':0}),justify="center"),
    ]
,className='app-page')

################################################################################

############################## ARITHMETICS #####################################
arithmeticsMenu = html.Div([
    dbc.Row([dbc.Col(html.H2(children='Arithmetics Benchmarks'))]),
    # Curves
    dbc.Row(
        [
            dbc.Col(html.H4(style={'text-align': 'center'}, children='Select Curve(s):'),
                xs={'size':'auto', 'offset':0}, sm={'size':'auto', 'offset':0}, md={'size':'auto', 'offset':3},
                lg={'size':'auto', 'offset':0}, xl={'size':'auto', 'offset':0}),
            dbc.Col(dcc.Dropdown(
                #style = {'text-align': 'center', 'font-size': '18px', 'width': '210px'},
                id='arithmetics-curves',
                options=arithmetics_curves,
                value=arithmetics_curves,
                clearable=False,
                multi=True),
                xs={'size':'auto', 'offset':0}, sm={'size':'auto', 'offset':0}, md={'size':'auto', 'offset':0},
                lg={'size':'auto', 'offset':0}, xl={'size':'auto', 'offset':0}),
        ],
    ),
    dbc.Row(dbc.Col(html.P(style={'font-size': '16px', 'opacity': '70%'},
        children='''Select which curves to display'''
    ))),
    # Frameworks
    dbc.Row(
        [
            dbc.Col(html.H4(style={'text-align': 'center'}, children='Select Framework(s):'),
                xs={'size':'auto', 'offset':0}, sm={'size':'auto', 'offset':0}, md={'size':'auto', 'offset':3},
                lg={'size':'auto', 'offset':0}, xl={'size':'auto', 'offset':0}),
            dbc.Col(dcc.Dropdown(
                #style = {'text-align': 'center', 'font-size': '18px', 'width': '210px'},
                id='arithmetics-frameworks',
                options=arithmetics_frameworks,
                value=arithmetics_frameworks,
                clearable=False,
                multi=True),
                xs={'size':'auto', 'offset':0}, sm={'size':'auto', 'offset':0}, md={'size':'auto', 'offset':0},
                lg={'size':'auto', 'offset':0}, xl={'size':'auto', 'offset':0}),
        ],
    ),
    dbc.Row(dbc.Col(html.P(style={'font-size': '16px', 'opacity': '70%'},
        children='''Select which frameworks to display'''
    ))),
    # Fields
    dbc.Row(
        [
            dbc.Col(html.H4(style={'text-align': 'center'}, children='Select Field(s):'),
                xs={'size':'auto', 'offset':0}, sm={'size':'auto', 'offset':0}, md={'size':'auto', 'offset':3},
                lg={'size':'auto', 'offset':0}, xl={'size':'auto', 'offset':0}),
            dbc.Col(dcc.Dropdown(
                #style = {'text-align': 'center', 'font-size': '18px', 'width': '210px'},
                id='arithmetics-fields',
                options=arithmetics_fields,
                value=arithmetics_fields,
                clearable=False,
                multi=True),
                xs={'size':'auto', 'offset':0}, sm={'size':'auto', 'offset':0}, md={'size':'auto', 'offset':0},
                lg={'size':'auto', 'offset':0}, xl={'size':'auto', 'offset':0}),
        ],
    ),
    dbc.Row(dbc.Col(html.P(style={'font-size': '16px', 'opacity': '70%'},
        children='''Select which field to display'''
    ))),
    dbc.Row(
        [
            dbc.Col(html.H4(style={'text-align': 'center'}, children='Select Operation:'),
                xs={'size':'auto', 'offset':0}, sm={'size':'auto', 'offset':0}, md={'size':'auto', 'offset':3},
                lg={'size':'auto', 'offset':0}, xl={'size':'auto', 'offset':0}),
            dbc.Col(dcc.Dropdown(
                style = {'text-align': 'center', 'font-size': '18px', 'width': '180px'},
                id='arithmetics-operation',
                options=arithmetics_operations,
                value=arithmetics_default_operation,
                clearable=False),
                xs={'size':'auto', 'offset':0}, sm={'size':'auto', 'offset':0}, md={'size':'auto', 'offset':0},
                lg={'size':'auto', 'offset':0}, xl={'size':'auto', 'offset':0}),

            dbc.Col(html.H4(style={'text-align': 'center', 'justify-self': 'right'}, children='Select Metric:'),
                xs={'size':'auto', 'offset':0}, sm={'size':'auto', 'offset':0}, md={'size':'auto', 'offset':3},
                lg={'size':'auto', 'offset':0}, xl={'size':'auto', 'offset':1}),
            dbc.Col(dcc.Dropdown(
                style = {'text-align': 'center', 'font-size': '18px', 'width': '180px'},
                id='arithmetics-metric',
                options=arithmetics_metrics,
                value=arithmetics_default_metric,
                clearable=False),
                xs={'size':'auto', 'offset':0}, sm={'size':'auto', 'offset':0}, md={'size':'auto', 'offset':0},
                lg={'size':'auto', 'offset':0}, xl={'size':'auto', 'offset':0}),
        ],
    ),
    html.Br(),
    dbc.Row(
        [
            dbc.Col(html.H4(style={'text-align': 'center'}, children='Select Input:'),
                xs={'size':'auto', 'offset':0}, sm={'size':'auto', 'offset':0}, md={'size':'auto', 'offset':3},
                lg={'size':'auto', 'offset':0}, xl={'size':'auto', 'offset':0}),
            dbc.Col(dcc.Dropdown(
                style = {'text-align': 'center', 'font-size': '18px', 'width': '500px'},
                id='arithmetics-input-dropdown',
                clearable=False),
                xs={'size':'auto', 'offset':0}, sm={'size':'auto', 'offset':0}, md={'size':'auto', 'offset':0},
                lg={'size':'auto', 'offset':0}, xl={'size':'auto', 'offset':0})
    ])
], className='menu')

arithmeticsLayout = html.Div([
    # Bar Chart of Benchmarks
    dbc.Row(dbc.Col(html.Div(id='arithmetics-bar'), xs={'size':'auto', 'offset':0}, sm={'size':'auto', 'offset':0}, md={'size':7, 'offset':0}, lg={'size':'auto', 'offset':0},
            xl={'size':10, 'offset':0}),justify="center"),
    # Line Chart
    dbc.Row(dbc.Col(html.Div(id='arithmetics-line'), xs={'size':'auto', 'offset':0}, sm={'size':'auto', 'offset':0}, md={'size':7, 'offset':0}, lg={'size':'auto', 'offset':0},
            xl={'size':10, 'offset':0}),justify="center"),
    ]
,className='app-page')

################################################################################
