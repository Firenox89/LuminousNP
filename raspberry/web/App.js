import React from 'react';
import {makeStyles} from '@material-ui/core/styles';

import {Button, Container} from '@material-ui/core';
import Card from '@material-ui/core/Card';
import CardContent from '@material-ui/core/CardContent';
import Typography from '@material-ui/core/Typography';
import Checkbox from "@material-ui/core/Checkbox";

import FormControlLabel from "@material-ui/core/FormControlLabel";
import FormControl from "@material-ui/core/FormControl";
import InputLabel from "@material-ui/core/InputLabel";
import Select from "@material-ui/core/Select";
import MenuItem from "@material-ui/core/MenuItem";

import {HuePicker} from 'react-color'

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
    formControl: {
        margin: theme.spacing(1),
        minWidth: 120,
    },
}));

export default function App() {
    const classes = useStyles();

    const [color, setColor] = React.useState('#fff');
    const [ledsOn, setLedsOn] = React.useState(true);
    const [useWhite, setUseWhite] = React.useState(true);
    const [effect, setEffect] = React.useState(0);

    const buildConfig = () => {
        return {
            config: {
                power: ledsOn,
                useWhite: useWhite,
                color: color,
                effect: effect
            }
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
            <Card className={classes.root}>
                <CardContent>
                    <Typography className={classes.title} gutterBottom>
                        Config
                    </Typography>
                    <Typography className={classes.subtitle} gutterBottom>
                        LEDs
                    </Typography>
                    <Button variant="contained" className={classes.button}
                            onClick={() => setLedsOn(false)}>Off</Button>
                    <Button variant="contained" className={classes.button} color="primary"
                            onClick={() => setLedsOn(true)}> On </Button>
                    <FormControlLabel
                        control={
                            <Checkbox
                                checked={useWhite}
                                onChange={(event => {
                                    setUseWhite(event.target.checked);
                                })}
                                value="checkedB"
                                color="primary"
                            />
                        }
                        label="Use white LED"
                    />
                    <Typography className={classes.subtitle} gutterBottom>
                        Effect
                    </Typography>
                    <FormControl variant="outlined" className={classes.formControl}>
                        <InputLabel id="demo-simple-select-outlined-label">
                            Effect
                        </InputLabel>
                        <Select
                            id="demo-simple-select-outlined"
                            value={effect}
                            onChange={(event => {
                                setEffect(event.target.value)
                            })}
                        >
                            <MenuItem value={0}>Fill</MenuItem>
                            <MenuItem value={1}>FadeInOut</MenuItem>
                            <MenuItem value={2}>RainbowFade</MenuItem>
                            <MenuItem value={3}>Rainbow</MenuItem>
                            <MenuItem value={4}>RunningRainbow</MenuItem>
                            <MenuItem value={5}>RunningColor</MenuItem>
                        </Select>
                    </FormControl>
                    <Typography className={classes.subtitle} gutterBottom>
                        Color
                    </Typography>
                    <HuePicker
                        color={color}
                        onChange={color => setColor(color.hex)}
                        onChangeComplete={color => setColor(color.hex)}
                    />
                </CardContent>
            </Card>
            <Button variant="contained" className={classes.button} onClick={onApply}>Apply</Button>

        </Container>
    );
}
