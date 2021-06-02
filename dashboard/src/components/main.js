import { Switch, Route } from 'react-router-dom';

import About from './about'
import Heatmap from './heatmap'
import Stats from './stats'
import LiveLogs from './live-logs'

const Main = () => {
    return (
        <Switch>
            <Route exact path='/' component={About}></Route>
            <Route exact path='/heatmap' component={Heatmap}></Route>
            <Route exact path='/stats' component={Stats}></Route>
            <Route exact path='/live-logs' component={LiveLogs}></Route>
        </Switch>
    );
};

export default Main;
