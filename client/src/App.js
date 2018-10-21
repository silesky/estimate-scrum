import React from 'react'
import { BrowserRouter as Router, Route, Switch } from "react-router-dom";

import Home from './Home'
import Poker from './Poker'

export default () => (
  <Router>
    <Switch>
      <Route exact path="/" component={Home} />
      <Route path="/poker" component={Poker} />
      <Route component={() => <h1>Sorry. 404!</h1>} />
     </Switch>
  </Router>
);
