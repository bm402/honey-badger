import React from 'react';
import { Switch, Route } from 'react-router-dom';

import About from './about'
import Heatmap from './heatmap'

const Main = () => {
    return (
        <Switch>
            <Route exact path='/' component={About}></Route>
            <Route exact path='/heatmap' component={Heatmap}></Route>
        </Switch>
    );
};

export default Main;
