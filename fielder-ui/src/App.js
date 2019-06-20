import React from "react";
import "./App.css";
import Dashboard from "./materialEx/page-layout-examples/dashboard/Dashboard";
import SignIn from "./materialEx/page-layout-examples/sign-in/SignIn";
import SignUp from "./materialEx/page-layout-examples/sign-up/SignUp";
import { Provider as ReduxProvider } from "react-redux";
import configureStore from "./modules/store";

const reduxStore = configureStore(window.REDUX_INITIAL_DATA);

function App() {
    return (
        <ReduxProvider store={reduxStore}>
            <div className="App">
                {/* <header className="App-header">
                    <img src={logo} className="App-logo" alt="logo" />
                    <p>
                        Edit <code>src/App.js</code> and save to reload.
                    </p>
                    <a
                        className="App-link"
                        href="https://reactjs.org"
                        target="_blank"
                        rel="noopener noreferrer"
                    >
                        Learn React
                    </a>
                </header> */}
            </div>
            <Dashboard />
            <SignIn />
            <SignUp />
        </ReduxProvider>
    );
}

export default App;
