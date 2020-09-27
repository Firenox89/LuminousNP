import React from 'react';
import {makeStyles} from '@material-ui/core/styles';
import {Container} from '@material-ui/core';
import {LEDConfig} from './components/LEDConfig'
import Grid from "@material-ui/core/Grid";
import {NodeList} from "./components/NodeList";

const useStyles = makeStyles(theme => ({
    root: {
        width: 400,
        backgroundColor: "black"
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
        margin: 8,
        backgroundColor: "white"
    },
    paper: {
        padding: 8,
        marginBottom: theme.spacing(2),
        backgroundColor: "grey"
    },
    formControl: {
        margin: theme.spacing(1),
        minWidth: 120,
    },
}));

export default function App() {
    const classes = useStyles();

    const [effectId, setEffectID] = React.useState(0);
    const [colorPaletteId, setColorPaletteID] = React.useState(0);
    const [color, setColor] = React.useState('#ff0000');
    const [brightness, setBrightness] = React.useState(100);
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
    const buildConfig = (power, effectId, colorPaletteId, color, brightness) => {
        return {
            config: {
                power: power,
                useWhite: true,
                color: color,
                effect: effectId,
                colorPaletteId: colorPaletteId,
                brightness: brightness
            },
            nodes: connectedDevices
                .filter((value) => selected.includes(value.ID))
                .map((value) => {
                    return {ID: value.ID}
                })
        };
    }

    const onApply = (power, effectId, colorPaletteId, color, brightness) => {
        const config = buildConfig(power, effectId, colorPaletteId, color, brightness)
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
                <Grid item xs={4} sm={4}>
                    <NodeList
                        classes={classes}
                        selected={selected}
                        setSelected={setSelected}
                        connectedDevices={connectedDevices}
                    />
                </Grid>
                <Grid item xs={8} sm={8}>
                    <LEDConfig classes={classes}
                               onOff={() => {
                                   onApply(false, 0, 0, color, brightness)
                               }}
                               selectedEffect={effectId}
                               onEffectChange={effectId => {
                                   setEffectID(effectId)
                                   onApply(true, effectId, colorPaletteId, color, brightness)
                               }}
                               onColorPaletteChange={colorPaletteId => {
                                   setColorPaletteID(colorPaletteId)
                                   onApply(true, effectId, colorPaletteId, color, brightness)
                               }}
                               color={color}
                               onColorChange={color => {
                                   setColor(color.hex)
                                   onApply(true, effectId, colorPaletteId, color.hex, brightness)
                               }}
                               brightness={brightness}
                               onChangeBrightness={brightness => {
                                   setBrightness(brightness)
                                   onApply(true, effectId, colorPaletteId, color, brightness)
                               }}/>
                </Grid>
            </Grid>
        </Container>
    );
}
