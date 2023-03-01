# Import Bootstrap from Dash
import dash_bootstrap_components as dbc
from dash import Input, Output, State, html, dcc
from app import app

SMALL_LOGO = "/assets/img/logo.png"

# Navigation Bar fucntion
def Navbar():
    navbar = dbc.Navbar(
         dbc.Container(
            [
                dbc.NavbarToggler(id="navbar-toggler", className="navbar-toggler-custom"),
                html.A(
                    dbc.Row(
                        [
                            dbc.Col(html.Img(src=SMALL_LOGO, height="50px")),
                            dbc.Col(dbc.NavbarBrand("", className="ms-2")),
                        ],
                        align="center",
                    ),
                    href="/",
                    style={"textDecoration": "none"},
                    className="navbar-brand",
                ),
                dbc.Collapse(
                    dbc.Nav(
                        [
                            dbc.NavItem(dbc.NavLink("Circuit Benchmarks", href='/circuit', style={"color": "#003262"})),
                            dbc.NavItem(dbc.NavLink("Arithmetic Benchmarks", href='/arithmetic', style={"color": "#003262"})),
                            dbc.NavItem(dbc.NavLink("Elliptic Curve Benchmarks", href='/ec', style={"color": "#003262"})),
                        ],
                        navbar=True,
                        className="ml-auto",
                        fill=True,
                        pills=True,  
                    ),
                    id="navbar-collapse",
                    navbar=True,
                    is_open=False,
                    className="justify-content-end",
                ),
            ]
        ),
        color="#FFFFFF",
        dark=True,
        sticky="top",
        className="mb-5",
        expand='lg',
    )
    return navbar

@app.callback(
    Output("navbar-collapse", "is_open"),
    [Input("navbar-toggler", "n_clicks")],
    [State("navbar-collapse", "is_open")],
)
def toggle_navbar_collapse(n, is_open):
    if n:
        return not is_open
    return is_open
