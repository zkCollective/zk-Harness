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
circuits_operations = list(set(circuits_df['operation']))
circuits_operations_options = [
        {
            "label": html.P([op], style={'font-size': 20, 'margin-left': '10px', 'margin-right': '10px', 'margin-bottom': '.5rem', 'display': 'inline-block'}),
            "value": op,
        }
    for op in circuits_operations
]
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

ec_curves = list(set(ec_df['curve']))
ec_frameworks = list(set(ec_df['framework']))
ec_operations = list(set(ec_df['operation']))
ec_default_operation = "pairing"
ec_metrics = ["time", "ram"]
ec_default_metric = "time"


def get_header(text):
    return html.Div(
    dbc.Row(
        [
            dbc.Col(
                html.H2(children=text),
                className="subsection-header text-6xl md:text-8xl underline text-center"
            )
        ]
    )
)
# get_curves(curves, 'circuits-curves')
def get_curves(curves, _id):
    return html.Div([
        # Curves
        dbc.Row([
            dbc.Col(
                html.H4(
                    style={'text-align': 'center'},
                    children='Select Curve(s):',
                    className='select-curve'
                ),
                xs={'size': 'auto', 'offset': 0},
                sm={'size': 'auto', 'offset': 0},
                md={'size': 'auto', 'offset': 3},
                lg={'size': 'auto', 'offset': 0},
                xl={'size': 'auto', 'offset': 0}
            ),
            dbc.Col(
                dcc.Dropdown(
                    id=_id,
                    options=curves,
                    value=curves,
                    clearable=False,
                    multi=True,
                    className='multi-option-style'
                ),
                xs={'size': 'auto', 'offset': 0},
                sm={'size': 'auto', 'offset': 0},
                md={'size': 'auto', 'offset': 0},
                lg={'size': True, 'offset': 0},
                xl={'size': True, 'offset': 0}
            ),
        ]),
        dbc.Row(
            dbc.Col(
                html.P(
                    style={'font-size': '16px', 'opacity': '70%'},
                    children='Select which curves to display'
                )
            )
        )
    ])


# get_frameworks(circuits_frameworks, 'circuits-frameworks')
def get_frameworks(frameworks, _id):
    return html.Div([
        dbc.Row([
            dbc.Col(
                html.H4(
                    style={'text-align': 'center'},
                    children='Select Framework(s):',
                    className='select-curve'
                ),
                xs={'size': 'auto', 'offset': 0},
                sm={'size': 'auto', 'offset': 0},
                md={'size': 'auto', 'offset': 3},
                lg={'size': 'auto', 'offset': 0},
                xl={'size': 'auto', 'offset': 0}
            ),
            dbc.Col(
                dcc.Dropdown(
                    id=_id,
                    options=frameworks,
                    value=frameworks,
                    clearable=False,
                    multi=True,
                    className='multi-option-style'
                ),
                xs={'size': 'auto', 'offset': 0},
                sm={'size': 'auto', 'offset': 0},
                md={'size': 'auto', 'offset': 0},
                lg={'size': True, 'offset': 0},
                xl={'size': True, 'offset': 0}
            ),
        ]),
        dbc.Row(
            dbc.Col(
                html.P(
                    style={'font-size': '16px', 'opacity': '70%'},
                    children='Select which frameworks to display'
                )
            )
        )
    ])


def get_backends_fields(backends, _id, case):
    return html.Div([
        dbc.Row([
            dbc.Col(
                html.H4(
                    style={'text-align': 'center'},
                    children='Select ' + case + ' :',
                    className='select-curve'
                ),
                xs={'size': 'auto', 'offset': 0},
                sm={'size': 'auto', 'offset': 0},
                md={'size': 'auto', 'offset': 3},
                lg={'size': 'auto', 'offset': 0},
                xl={'size': 'auto', 'offset': 0}
            ),
            dbc.Col(
                dcc.Dropdown(
                    id=_id,
                    options=backends,
                    value=backends,
                    clearable=False,
                    multi=True,
                    className='multi-option-style'
                ),
                xs={'size': 'auto', 'offset': 0},
                sm={'size': 'auto', 'offset': 0},
                md={'size': 'auto', 'offset': 0},
                lg={'size': True, 'offset': 0},
                xl={'size': True, 'offset': 0}
            ),
        ]),
        dbc.Row(
            dbc.Col(
                html.P(
                    style={'font-size': '16px', 'opacity': '70%'},
                    children='Select which ' + str(case) + ' to display'
                )
            )
        )
    ])


def get_operation(operation, _id_operation, case):
    if _id_operation == 'circuits-circuit':
        default = circuits_default_circuit
    elif _id_operation == 'arithmetics-operation':
        default = arithmetics_default_operation     
    elif _id_operation == 'ec-operation':
        default = ec_default_operation
    
    return html.Div([
        dbc.Row([
            dbc.Col(
                html.H4(
                    style={'text-align': 'center'},
                    children=f'Select {case}:',
                    className='select-curve'
                ),
                xs={'size': 'auto', 'offset': 0},
                sm={'size': 'auto', 'offset': 0},
                md={'size': 'auto', 'offset': 3},
                lg={'size': 'auto', 'offset': 0},
                xl={'size': 'auto', 'offset': 0}
            ),
            dbc.Col(
                dcc.Dropdown(
                    id=_id_operation,
                    options=operation,
                    value=default,
                    clearable=False,
                    className='dropdown-class-input'
                ),
                xs={'size': True, 'offset': 0},
                sm={'size': True, 'offset': 0},
                md={'size': True, 'offset': 0},
                lg={'size': True, 'offset': 0},
                xl={'size': True, 'offset': 0}
            ),
        ])
    ])

def get_metric(metric, _id_metric, case):
    if _id_metric == 'circuits-metric':
        default = circuits_default_metric
    elif _id_metric == 'arithmetics-metric':
        default = arithmetics_default_metric
    elif _id_metric == 'ec-metric':
        default = ec_default_metric
    
    return html.Div([
        dbc.Row([
            dbc.Col(
                html.H4(
                    style={'text-align': 'left'},
                    children='Select Metric:',
                    className='select-curve'
                ),
                xs={'size': 'auto', 'offset': 0},
                sm={'size': 'auto', 'offset': 0},
                md={'size': 'auto', 'offset': 3},
                lg={'size': 'auto', 'offset': 0},
                xl={'size': 'auto', 'offset': 0}
            ),
            dbc.Col(
                dcc.Dropdown(
                    id=_id_metric,
                    options=metric,
                    value=default,
                    clearable=False,
                    className='dropdown-class-input'
                ),
                xs={'size': True, 'offset': 0},
                sm={'size': True, 'offset': 0},
                md={'size': True, 'offset': 0},
                lg={'size': True, 'offset': 0},
                xl={'size': True, 'offset': 0}
            )
        ])
    ])

# get_curve_operation_checks(circuits_operations_options, circuits_operations, 'circuits-operation')
def get_curve_operation_checks(operations, options, _id):
    return html.Div([
        dbc.Row([
            dbc.Col(html.H4(style={'text-align': 'center'}, 
                            children='Select Operation(s):',
                            className='select-curve'
                    ),
                    xs={'size': 'auto', 'offset': 0},
                    sm={'size': 'auto', 'offset': 0},
                    md={'size': 'auto', 'offset': 3},
                    lg={'size': 'auto', 'offset': 0},
                    xl={'size': 'auto', 'offset': 0}),
            dbc.Col(dcc.Checklist(
                    operations,
                    id=_id,
                    value=options,
                    inline=True,
                    labelStyle={
                        'display': 'block',
                        'margin-right': '1rem',
                        'padding': '0.5rem 0',
                        'color': '#003262',
                        'font-size': '12px'
                    },
                    inputStyle={
                        'margin-right': '0.5rem'
                    }
                ),
                xs={'size': 'auto', 'offset': 0},
                sm={'size': 'auto', 'offset': 0},
                md={'size': 'auto', 'offset': 0},
                lg={'size': 'auto', 'offset': 0},
                xl={'size': 'auto', 'offset': 0})
        ])
    ])


def get_input(_id):
    return html.Div([
        dbc.Row([
            dbc.Col(
                html.H4(
                    style={'text-align': 'center'},
                    children='Select Input:',
                    className='select-curve'
                ),
                xs={'size': 'auto', 'offset': 0},
                sm={'size': 'auto', 'offset': 0},
                md={'size': 'auto', 'offset': 3},
                lg={'size': 'auto', 'offset': 0},
                xl={'size': 'auto', 'offset': 0}
            )
        ]),
        dbc.Row([
            dbc.Col(
                dcc.Dropdown(
                    id=_id,
                    clearable=False,
                    className='dropdown-class-input'
                ),
                xs={'size': True, 'offset': 0},
                sm={'size': True, 'offset': 0},
                md={'size': True, 'offset': 0},
                lg={'size': True, 'offset': 0},
                xl={'size': True, 'offset': 0}
            )
        ])
    ])


def get_count(_id):
    return html.Div([
        dbc.Row([
            dbc.Col(
                html.H4(
                    style={'text-align': 'center'},
                    children='Number of benchmarks executed: ',
                    className='select-curve'
                ),
                xs={'size': 'auto', 'offset': 0},
                sm={'size': 'auto', 'offset': 0},
                md={'size': 'auto', 'offset': 3},
                lg={'size': 'auto', 'offset': 0},
                xl={'size': 'auto', 'offset': 0}
            ),
            dbc.Col(
                html.H5(
                    id=_id,
                    className='select-curve'
                ),
                xs={'size': 'auto', 'offset': 0},
                sm={'size': 'auto', 'offset': 0},
                md={'size': 'auto', 'offset': 0},
                lg={'size': 'auto', 'offset': 0},
                xl={'size': 'auto', 'offset': 0}
            )
        ])
    ])


################################# CIRCUITS #####################################
circuitMenu = html.Div([
    get_header('Circuit Benchmarks'),
    get_curves(circuits_curves, 'circuits-curves'),
    get_frameworks(circuits_frameworks, 'circuits-frameworks'),
    get_backends_fields(circuits_backends, 'circuits-backends', 'backend'),
    get_operation(circuits_circuits, 'circuits-circuit','circuit'),
    html.Br(),
    get_metric(circuits_metrics, 'circuits-metric', 'circuit'),
    html.Br(),
    # Curve operations checkboxes
    get_curve_operation_checks(circuits_operations_options, circuits_operations, 'circuits-operation'),
    html.Br(),
    get_input('circuits-input-dropdown')
], className='menu')

circuitLayout = html.Div([
    # Bar Chart of Benchmarks
    dbc.Row(dbc.Col(html.Div(id='circuits-bar'), xs={'size':'auto', 'offset':0}, sm={'size':'auto', 'offset':0}, md={'size':7, 'offset':0}, lg={'size':'auto', 'offset':0},
            xl={'size':10, 'offset':0}),justify="center"),
    # Line Chart
    dbc.Row(dbc.Col(html.Div(id='circuits-line'), xs={'size':'auto', 'offset':0}, sm={'size':'auto', 'offset':0}, md={'size':7, 'offset':0}, lg={'size':'auto', 'offset':0},
            xl={'size':10, 'offset':0}),justify="center"),
    # Cicrcuit Constraint Table
    dbc.Row(dbc.Col(html.Div(id='circuits-data'), xs={'size':'auto', 'offset':0}, sm={'size':'auto', 'offset':0}, md={'size':7, 'offset':0}, lg={'size':'auto', 'offset':0},
            xl={'size':10, 'offset':0}),justify="center"),
    ]
    ,className='app-page'
)

################################################################################

############################## ARITHMETICS #####################################
arithmeticsMenu = html.Div([
    get_header('Arithmetics Benchmarks'),
    get_curves(arithmetics_curves, 'arithmetics-curves'),
    get_frameworks(arithmetics_frameworks, 'arithmetics-frameworks'),
    get_backends_fields(arithmetics_fields, 'arithmetics-fields', 'field'),
    get_operation(arithmetics_operations, 'arithmetics-operation', 'operation'),
    html.Br(),
    get_metric(arithmetics_metrics, 'arithmetics-metric', 'operation'),
    html.Br(),
    get_input('arithmetics-input-dropdown'),
    html.Br(),
    get_count('arithmetics-count'),
    ]
    , className='menu'
)

arithmeticsLayout = html.Div([
    # Bar Chart of Benchmarks
    dbc.Row(dbc.Col(html.Div(id='arithmetics-bar'), xs={'size':'auto', 'offset':0}, sm={'size':'auto', 'offset':0}, md={'size':7, 'offset':0}, lg={'size':'auto', 'offset':0},
            xl={'size':10, 'offset':0}),justify="center"),
    ]
    ,className='app-page'
)

################################################################################

################################### EC #########################################
ecMenu = html.Div([
    get_header('Elliptic Curves Benchmarks'),
    get_curves(ec_curves, 'ec-curves'),
    get_frameworks(ec_frameworks, 'ec-frameworks'),
    get_operation(ec_operations, 'ec-operation', 'operation'),
    html.Br(),
    get_metric(ec_metrics, 'ec-metric', 'operation'),
    html.Br(),
    get_input('ec-input-dropdown'),
    html.Br(),
    get_count('ec-count'),
    ]
    , className='menu'
)

ecLayout = html.Div([
    # Bar Chart of Benchmarks
    dbc.Row(dbc.Col(html.Div(id='ec-bar'), xs={'size':'auto', 'offset':0}, sm={'size':'auto', 'offset':0}, md={'size':7, 'offset':0}, lg={'size':'auto', 'offset':0},
            xl={'size':10, 'offset':0}),justify="center"),
    ]
    ,className='app-page'
)

################################################################################

