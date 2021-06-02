import Button from 'react-bootstrap/Button'

import './about.css'

const About = () => {
    return (
        <div className="about-page">
            <h2>Welcome to Honey Badger</h2>
            <p>
                The internet is a big mysterious black box of machines communicating with one another.
                Ever wondered how much communication is going on behind your back? Or which far corner
                of the world is getting in touch with your home router to see whether it is running
                outdated software and is vulnerable to a cyber attack?
            </p>
            <p>
                Honey Badger is a honeypot application that allows you to see just how much traffic an
                open server on the internet is subjected to. It spins up a random server on the AWS
                cloud and listens for any attempted connections made to that server. Every attempted
                connection is stored by the application and is used to keep track of who is trying to
                communicate with our server, where in the world the requests are coming from, and what
                payload data is being sent with the request.
            </p>
            <p>
                This dashboard shows collections of these metrics displayed in different ways. We have
                a heatmap which shows the distribution of connection attempts across the world map, a
                stats page which shows who is trying to connect to our server the most stubbornly, and
                a live logs page where you can see details of the connection attempts in real time.
            </p>
            <Button variant="dark" href="https://github.com/bncrypted/honey-badger">View project on GitHub</Button>
        </div>
    );
};

export default About;
