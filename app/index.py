# import dash-core, dash-html, dash io, bootstrap
from dash import dcc, html
from dash.dependencies import Input, Output

# Dash Bootstrap components
import dash_bootstrap_components as dbc

# Navbar, layouts, custom callbacks
from layouts import circuitMenu, circuitLayout, arithmeticsMenu, arithmeticsLayout, ecMenu, ecLayout
# Need to import it here so callbacks are loaded
import callbacks

from navbar import Navbar
# Import app
from app import app as application
# Import server for deployment
# from app import srv as server

app = application

# Layout variables, navbar, header, content, and container
nav = Navbar()

footer = dbc.Row(
    dbc.Col(
        html.Div([
                "View Source on ",
                html.A(
                    href="https://github.com/zkCollective/zk-Harness",
                    children=html.Img(src="https://github.githubassets.com/images/modules/logos_page/GitHub-Mark.png", style={"height": "1.5em"})
                ),
            ], style={"float": "center", "font-family": "Courier New, monospace" , 'color': '#003262'})
        ),className='banner')

content = html.Div([
    dcc.Location(id='url'),
    html.Div(id='page-content')
])

container = dbc.Container([
    content,
    footer,
])


# Menu callback, set and return
# Declair function  that connects other pages with content to container
@app.callback(Output('page-content', 'children'),
            [Input('url', 'pathname')])
def display_page(pathname):
    if pathname == '/':
        return html.Div([dcc.Markdown('''
            ### What is it?
            zk-Harness is a benchmarking framework for general purpose zero-knowledge proofs. 
            It is designed to be modular - new circuit implementations and ZKP-frameworks can be easily added, 
            without extensive developer overhead. zk-Harness has a standardized set of interfaces for configuring 
            benchmark jobs and formatting log outputs.

            __NOTE:__ zk-Harness is a WIP. We welcome and value contributions from all individuals. You can find our contribution guidelines on [GitHub](https://github.com/zkCollective/zk-Harness).

            ### Main Features


            zk-Harness currently includes the following:


            - Benchmarks for field arithmetic
            - Benchmarks for Elliptic curve group operations
            - Benchmarks for circuit implementations
            - In the following frameworks:
                - gnark
                - circom

            A detailed list of included sub-components and the framework architecture can be found in the [GitHub](https://github.com/zkCollective/zk-Harness) repository.

            ### FAQ
            **What data am I looking at?**

            The data you are looking at are measurements of common ZKP frameworks, currently executed on a specific local processor.
            Our source-code is open-source and you can find the raw data [here](https://github.com/zkCollective/zk-Harness/tree/main/benchmarks). 

            **Which criterias is used to determine whether a ZKP framework is included?**

            We do not favor any proving system over another, we aim for completenes such that developers may benefit from a leveled and standardized comparison.

            **How do you determine the validity of the data?**

            Before including a new set of benchmarks, we ensure the compliance of the benchmarks with our standardized interfaces and specification of metrics.
            Further, we manually check the sanity of the measurements, such that equivalent measures are applied for comparability.
            In some cases, such as for Circom and SnarkJS, the comparison is not equal due to the overhead of certain components (e.g., see [this](https://github.com/zkCollective/zk-Harness/issues/1) issue). 

            Of course, there is always the possibility of a bug. If you find something suspicious, please let us know by opening an issue in our [GitHub](https://github.com/zkCollective/zk-Harness).

            **How can I contribute?**

            zk-Harness is an open-source public good developed as initiative by the [zk-Collective](https://zkcollective.org/). 
            Currently, zk-Harness is a part of the [ZKP/Web3 Hackathon](https://rdi.berkeley.edu/zkp-web3-hackathon/) - you can find dedicated tasks with specific prices [here](https://drive.google.com/file/d/1Igm47dFXSOFAC_wldfUG4Y9OiITqlbQu/view?usp=share_link).
            If you'd like to make a contribution by including a new system, please see the documentation on [How to contribute?](https://github.com/zkCollective/zk-Harness).
            
        ''')],className='home', style={'text-align': 'justify', 'font-size': '14px', 'color': '#003262'})
    elif pathname == '/circuit':
        return circuitMenu, circuitLayout
    elif pathname == '/arithmetic':
        return arithmeticsMenu, arithmeticsLayout
    elif pathname == '/ec':
        return ecMenu, ecLayout
    else:
        # If the user tries to reach a different page, return a 404 message
        return html.Div(
            [
                html.H1("404: Not found", className="text-danger"),
                html.Hr(),
                html.P(f"The pathname {pathname} was not recognised..."),
            ]
        )


# Main index function that will call and return all layout variables
def index():
    layout = html.Div([
            nav,
            container
        ])
    return layout

# Set layout to index function
app.layout = index()

# Call app server
if __name__ == '__main__':
    # set debug to false when deploying app
    srv = app.server
    app.run_server(debug=True)
else:
    srv = app.server
