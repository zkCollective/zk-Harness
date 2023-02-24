# Import Bootstrap from Dash
import dash_bootstrap_components as dbc


# Navigation Bar fucntion
def Navbar():
    navbar = dbc.NavbarSimple(children=[
            dbc.NavItem(dbc.NavLink("Circuit Benchmarks", href='/circuit')),
            dbc.NavItem(dbc.NavLink("Arithmetic Benchmarks", href='/arithmetic')),
            dbc.NavItem(dbc.NavLink("Elliptic Curve Benchmarks", href='/ec')),
        ],
        brand="Home",
        brand_href="/",
        sticky="top",
        color="light",
        dark=False,
        expand='lg',)
    return navbar
