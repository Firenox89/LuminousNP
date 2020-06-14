import React from 'react';
import {makeStyles} from '@material-ui/core/styles';
import {Container} from '@material-ui/core';
import {LEDConfig} from './components/LEDConfig'
import NodeTable from "./components/NodeTable";
import Grid from "@material-ui/core/Grid";

const useStyles = makeStyles(theme => ({
    root: {
        width: 400,
    },
    bullet: {
        display: 'inline-block',
        margin: '0 2px',
        transform: 'scale(0.8)',
    },
    title: {
        fontSize: 22,
    },
    subtitle: {
        fontSize: 18,
    },
    pos: {
        marginBottom: 12,
    },
    button: {
        margin: 8
    },
    paper: {
        width: '100%',
        padding: 8,
        marginBottom: theme.spacing(2),
    },
    formControl: {
        margin: theme.spacing(1),
        minWidth: 120,
    },
}));

export default function App() {
    const classes = useStyles();

    const [color, setColor] = React.useState('#ff0000');
    const [ledsOn, setLedsOn] = React.useState(true);
    const [useWhite, setUseWhite] = React.useState(true);
    const [effect, setEffect] = React.useState(0);

    const [connectedDevices, setConnectedDevices] = React.useState([]);

    const [selected, setSelected] = React.useState([]);

    React.useEffect(() => {
        const timeout = 1000
        const interval = setInterval(() => {
            fetch("/getConnectedNodeMCUs")
                .then(response => response.json())
                .then(data => {
                    setConnectedDevices(data)
                })
        }, timeout);
        return () => {
            clearInterval(interval);
        };
    }, []);

    //console.log("device", connectedDevices)
    const buildConfig = () => {
        return {
            config: {
                power: ledsOn,
                useWhite: useWhite,
                color: color.substring(1),//cut the #
                effect: effect
            },
            nodes: connectedDevices
                .filter((value) => selected.includes(value.ID))
                .map((value) => {
                    return {ID: value.ID}
                })
        };
    }

    const onApply = () => {
        const config = buildConfig()
        const configString = JSON.stringify(config)
        console.log(configString)

        fetch('/setConfig', {
            method: 'POST', // *GET, POST, PUT, DELETE, etc.
            mode: 'cors', // no-cors, *cors, same-origin
            cache: 'no-cache', // *default, no-cache, reload, force-cache, only-if-cached
            credentials: 'same-origin', // include, *same-origin, omit
            headers: {
                'Content-Type': 'application/json'
                // 'Content-Type': 'application/x-www-form-urlencoded',
            },
            redirect: 'follow', // manual, *follow, error
            referrerPolicy: 'no-referrer', // no-referrer, *no-referrer-when-downgrade, origin, origin-when-cross-origin, same-origin, strict-origin, strict-origin-when-cross-origin, unsafe-url
            body: configString // body data type must match "Content-Type" header
        }).then(data => {
            console.log(data); // JSON data parsed by `response.json()` call
        });
    };

    return (
        <Container>
            <Grid container spacing={1}>
                <Grid item xs={12} sm={8}>
                    <NodeTable
                        selected={selected}
                        setSelected={setSelected}
                        connectedDevices={connectedDevices}
                    />
                </Grid>
                <Grid item xs={12} sm={4}>
                    <LEDConfig classes={classes}
                               power={ledsOn}
                               setPower={event => {
                                   setLedsOn(event.target.checked);
                               }}
                               useWhite={useWhite}
                               setUseWhite={event => {
                                   setUseWhite(event.target.checked);
                               }}
                               selectedEffect={effect}
                               onEffectChange={event => {
                                   setEffect(event.target.value)
                               }}
                               color={color}
                               onColorChange={color => setColor(color.hex)}
                               onApply={onApply}/>
                </Grid>
            </Grid>
        </Container>
    );
}
