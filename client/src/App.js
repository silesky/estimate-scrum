import React from 'react'
import { BrowserRouter as Router, Route, Switch } from "react-router-dom";
import './App.css';
import Home from './Pages/Home'
import Session from './Pages/Session'

export default () => (
  <Router>
    <Switch>
      <Route exact path="/" component={Home} />
      <Route path="/session" component={Session} />
      <Route component={() => <h1>Sorry. 404!</h1>} />
     </Switch>
  </Router>
);
