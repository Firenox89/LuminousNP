import React from 'react';
import {makeStyles} from '@material-ui/core/styles';

import {Button, Container} from '@material-ui/core';
import {LEDConfig} from './LEDConfig'
import List from '@material-ui/core/List';
import ListItem from '@material-ui/core/ListItem';
import ListItemIcon from '@material-ui/core/ListItemIcon';
import ListItemText from '@material-ui/core/ListItemText';
import Checkbox from '@material-ui/core/Checkbox';
import Card from "@material-ui/core/Card";
import CardContent from "@material-ui/core/CardContent";
import Typography from "@material-ui/core/Typography";

const useStyles = makeStyles(theme => ({
    root: {
        width: 400,
        margin: 16,
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
    const [checked, setChecked] = React.useState([]);

    const handleToggle = (value) => () => {
        const currentIndex = checked.indexOf(value);
        const newChecked = [...checked];

        if (currentIndex === -1) {
            newChecked.push(value);
        } else {
            newChecked.splice(currentIndex, 1);
        }

        setChecked(newChecked);
    };
    if (connectedDevices.length === 0) {
        fetch("/getConnectedNodeMCUs")
            .then(response => response.json())
            .then(data => {
                console.log(data)
                setConnectedDevices(data)
            })
    }

    console.log("device", connectedDevices)
    const buildConfig = () => {
        return {
            config: {
                power: ledsOn,
                useWhite: useWhite,
                color: color.substring(1),//cut the #
                effect: effect
            },
            nodes: connectedDevices
                .filter((value, index) => checked.includes(index))
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
                       onColorChange={color => setColor(color.hex)}/>

            <Card className={classes.root}>
                <CardContent>
                    <Typography className={classes.title} gutterBottom>
                        Nodes
                    </Typography>
                    <List className={classes.root}>
                        {connectedDevices.map((value, index) => {
                            const labelId = `checkbox-list-label-${index}`;

                            return (
                                <ListItem key={index} role={undefined} dense button onClick={handleToggle(index)}>
                                    <ListItemIcon>
                                        <Checkbox
                                            edge="start"
                                            checked={checked.indexOf(index) !== -1}
                                            tabIndex={-1}
                                            disableRipple
                                            inputProps={{'aria-labelledby': labelId}}
                                        />
                                    </ListItemIcon>
                                    <ListItemText id={labelId} primary={`${value.ID}`}/>
                                    <ListItemText id={labelId} primary={`IP: ${value.IP}`}/>
                                    <ListItemText id={labelId} primary={`LEDS: ${value.LedCount}`}/>
                                </ListItem>
                            );
                        })}
                    </List>
                </CardContent>
            </Card>
            <Button variant="contained" className={classes.button} onClick={onApply}>Apply</Button>

        </Container>
    );
}
