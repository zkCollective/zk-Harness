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

ARCHITECTURE_PATH = "/assets/img/HarnessSpecification.png"

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
        return html.Div([
            dbc.Row(
                [
                    dbc.Col(
                        html.H2(children='zk-Harness'),
                        className="subsection-header text-6xl md:text-8xl underline text-center"
                    )
                ]
            ),
            dcc.Markdown('''
            
            We cordially invite the zk SNARK community to join us in creating a comprehensive benchmarking framework (zk-Harness) for zk SNARKs. This is a crucial step in the important mission to create a reference point for non-experts and experts alike on what zkSNARK scheme best suits their needs, and to also promote further research by identifying performance gaps. We believe that the collective efforts of the community will help to achieve this goal. Whether you are a researcher, developer, or simply passionate about zk SNARKs, we welcome your participation and contribution in this exciting initiative.

            It is designed to be modular - new circuit implementations and ZKP-frameworks can be easily added, without extensive developer overhead.
            zk-Harness has a standardized set of interfaces for configuring benchmark jobs and formatting log outputs.
            Once a new component is included, it's benchmarks will be displayed on this website.

            __NOTE:__ zk-Harness is a WIP. We welcome and value contributions from all individuals. You can find our contribution guidelines on [GitHub](https://github.com/zkCollective/zk-Harness).

            ### Main Features

            There is a large and ever-increasing number of SNARK implementations. Although the
            theoretical complexity of the underlying proof systems is well understood, the concrete costs
            rely on a number of factors such as the efficiency of the field and curve implementations, the
            underlying proof techniques, and the computation model and its compatibility with the
            specific application. To elicit the concrete performance differences in different proof systems,
            it is important to separately benchmark the following:

            #### Field and Curve computations
            All popular SNARKs operate over prime fields, which are basically integers modulo p,
            i.e,. F_p. While some SNARKs are associated with a single field F_p, there are many
            SNARKs that rely on elliptic curve groups for security. For such SNARKs, the scalar
            field of the elliptic curve is F_p, and the base field is a different field F_q. Thus, the
            aim is to benchmark the field F_p, along with the field F_q and the elliptic curve
            group (if applicable).
            Benchmarking F_p and F_q involves benchmarking the following operations:
           
            - Addition
            - Subtraction
            - Multiplication
            - (Modular) Exponentiation
            - Inverse Exponentiation

            An elliptic curve is defined over a prime field of specific order (F_q). The elliptic curve
            group (E(F_q)) consists of the subgroup of points in the field that are on the curve,
            including a special point at infinity. While some SNARKs operate over elliptic curves
            without requiring pairings, others require pairings and therefore demand for
            pairing-friendly elliptic curves. The pairing operation takes an element from G_1 and
            an element from G_2 and computes an element in G_T. The elements of G_T are
            typically denoted by e(P, Q), where P is an element of G_1 and Q is an element of
            G_2. For efficiency, it is required that not only is the finite field arithmetic fast, but also
            the arithmetic in groups G_1 and G_2 as well as pairings are efficient. Therefore, we
            intend to benchmark the following operations over pairing-friendly elliptic curves:


            - Scalar Multiplication
                - in G for single elliptic curves
                - in G_1 and G_2 for pairing-friendly elliptic curves
            - Multi-Scalar Multiplication (MSM)
                - in G for single elliptic curves
                - in G_1 and G_2 for pairing-friendly elliptic curves
            - Parings
                - for pairing-friendly elliptic curves
            
            #### Circuits

            Many end-to-end applications require proving a specific cryptographic primitive,
            which requires the specification of said cryptographic primitive in a specific ZKP
            framework.

            - *Circuits for native field operations* - 
            These operations, namely, addition and multiplication in F_p, are supported
            by each SNARK library, and they are the most efficient to prove with a
            SNARK because arithmetic modulo F_p is the native computation model of a
            SNARK. This provides a good understanding of the efficiency of the core
            SNARK implementation.
            - *Circuits for non-native field operations* - 
            All computations we want to prove do not belong to arithmetic modulo p. For
            instance, Z_{2^64} or uint64/int64 is a popular data type in traditional
            programming languages. Or, we might want to prove arithmetic on a different
            field, say Z_q. This usually happens when we want to verify elliptic-curve
            based cryptographic primitives. An example of this is supporting verification of
            ECDSA signatures. The native field of elliptic curve underlying the chosen
            SNARK typically differs from the base field of the secp256k1 curve
            - *Circuits for SNARK-optimized primitives* - One of the challenges in the practically using SNARKs is their inefficiency
            with regard to traditional hash algorithms, like SHA-2, and traditional
            signature algorithms, such as ECDSA. They are fast when executed on a
            CPU, but prohibitively slow when used in a SNARK. As a result, the
            community has proposed several hash functions and signature algorithms
            that are SNARK-friendly, such as the following:
                - Poseidon Hash
                - Pedersen Hash
                - MIMC Hash
                - Ed25519 (EdDSA signature)
            - *Circuits for CPU-optimized primitives* - 
            Even though it would be beneficial to only rely on SNARK optimized
            primitives, practical applications often don’t allow for the usage of e.g.
            Poseidon hash functions or SNARK friendly signature schemes. For example,
            verifying ECDSA signatures in SNARKs is crucial when building e.g.
            zkBridge, however an implementation requires for non-native field arithmetic
            (see here), and therefore yields many constraints. Similarly, for building
            applications such as TLS Notary, one has to prove SHA-256 hash functions
            and AES-128 encryption which yields many constraints. Hence, we intend to
            benchmark the performance of the following cryptographic primitives and their
            circuit implementations in different ZKP-frameworks:
                - SHA-256
                - Blake2
                - ECDSA

            '''),
            html.Img(src=ARCHITECTURE_PATH, style={
                'width': '80%',
                'text-align': 'center',
                'display': 'block',
                'margin': 'auto',
                'padding': '20px 0'}),
            dcc.Markdown('''
            #### Current Features

            On a high level, zk-Harness takes as input a configuration file. The “Config Reader” reads
            the standardized config and invokes the ZKP framework as specified in the configuration file.
            You can find a description of the configuration file in the tutorials/config sub-folder of the
            GitHub repository. Each integrated ZKP framework exposes a set of functions that take as
            an input the standardized configuration parameters to execute the corresponding
            benchmarks. The output of benchmarking a given ZKP framework is a log file in csv format
            with standardized metrics. The log file is read by the “Log Analyzer”, which compiles the logs
            into pandas dataframes that are used by the front-end and displayed on the public website.
            You can find the standardized logging format in the tutorials/logs sub-folder.
            
            Currently, zk-Harness includes the following components as a starting point:

            - Benchmarks for field arithmetic
            - Benchmarks for Elliptic curve group operations
            - Benchmarks for circuit implementations
            - In the following frameworks:
                - gnark
                - circom

            We aim to successively expand this list to further include benchmarks for other ZKP frameworks, recursion and zk-EVMs.
            As a part of the ZKP/Web3 Hackathon hosted by UC Berkeley RDI, we aim to further develop the frameworks integrated into zk-Harness.
            You can find the program description detailing future integrations [here](https://drive.google.com/file/d/1Igm47dFXSOFAC_wldfUG4Y9OiITqlbQu/view).
            A detailed list of currently included sub-components and the framework architecture can be found in the [GitHub](https://github.com/zkCollective/zk-Harness) repository.

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
            Currently, zk-Harness is a part of the [ZKP/Web3 Hackathon](https://rdi.berkeley.edu/zkp-web3-hackathon/) - you can find dedicated tasks with specific prizes [here](https://drive.google.com/file/d/1Igm47dFXSOFAC_wldfUG4Y9OiITqlbQu/view?usp=share_link).
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
