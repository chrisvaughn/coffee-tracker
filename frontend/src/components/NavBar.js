import React from "react";
import { useAuth0 } from "../react-auth0-spa";
import { Link } from "react-router-dom";

const NavBar = () => {
    const { isAuthenticated, loginWithRedirect, logout } = useAuth0();

    return (
        <div>
            {!isAuthenticated && (
                <button onClick={() => loginWithRedirect({})}>Log in</button>
            )}

            {isAuthenticated && <button onClick={() => logout()}>Log out</button>}

            {/* NEW - add a link to the home and profile pages */}
            {isAuthenticated && (
                <span>
                    &nbsp;
                    <Link to="/">Home</Link>&nbsp;
                    <Link to="/profile">Profile</Link>&nbsp;
                    <Link to="/external-api">External API</Link>&nbsp;
                </span>
            )}
        </div>
    );
};

export default NavBar;