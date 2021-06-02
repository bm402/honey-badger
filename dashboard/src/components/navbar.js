import BSNavbar from 'react-bootstrap/Navbar';
import BSNav from 'react-bootstrap/Nav'

const Navbar = () => {
    return (    
        <BSNavbar bg="light" expand="lg">
            <BSNavbar.Brand href="/honey-badger/#/">Honey Badger</BSNavbar.Brand>
            <BSNavbar.Toggle aria-controls="basic-navbar-nav" />
            <BSNavbar.Collapse id="basic-navbar-nav">
                <BSNav>
                    <BSNav.Link href="/honey-badger/#/heatmap">Heatmap</BSNav.Link>
                </BSNav>
                <BSNav>
                    <BSNav.Link href="/honey-badger/#/stats">Stats</BSNav.Link>
                </BSNav>
                <BSNav>
                    <BSNav.Link href="/honey-badger/#/live-logs">Live logs</BSNav.Link>
                </BSNav>
            </BSNavbar.Collapse>
        </BSNavbar>
    );
};

export default Navbar;
