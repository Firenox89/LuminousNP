import Card from "@material-ui/core/Card";
import CardContent from "@material-ui/core/CardContent";
import Typography from "@material-ui/core/Typography";
import {Button} from "@material-ui/core";
import FormControlLabel from "@material-ui/core/FormControlLabel";
import Checkbox from "@material-ui/core/Checkbox";
import FormControl from "@material-ui/core/FormControl";
import InputLabel from "@material-ui/core/InputLabel";
import Select from "@material-ui/core/Select";
import MenuItem from "@material-ui/core/MenuItem";
import {HuePicker, SliderPicker} from "react-color";
import * as PropTypes from "prop-types";
import React from "react";

export function LEDConfig(props) {
    return <Card className={props.classes.root}>
        <CardContent>
            <Typography className={props.classes.title} gutterBottom>
                Config
            </Typography>
            <Typography className={props.classes.subtitle} gutterBottom>
                LEDs
            </Typography>
            <FormControlLabel
                control={
                    <Checkbox
                        checked={props.power}
                        onChange={(props.setPower)}
                        value="checkedB"
                        color="primary"
                    />
                }
                label="Power"
            />
            <FormControlLabel
                control={
                    <Checkbox
                        checked={props.useWhite}
                        onChange={(props.setUseWhite)}
                        value="checkedB"
                        color="primary"
                    />
                }
                label="Use white LED"
            />
            <Typography className={props.classes.subtitle} gutterBottom>
                Effect
            </Typography>
            <FormControl variant="outlined" className={props.classes.formControl}>
                <InputLabel id="demo-simple-select-outlined-label">
                    Effect
                </InputLabel>
                <Select
                    id="demo-simple-select-outlined"
                    value={props.selectedEffect}
                    onChange={(props.onEffectChange)}
                >
                    <MenuItem value={0}>Just White</MenuItem>
                    <MenuItem value={1}>Fill</MenuItem>
                    <MenuItem value={2}>FadeInOut</MenuItem>
                    <MenuItem value={3}>RainbowFade</MenuItem>
                    <MenuItem value={4}>Rainbow</MenuItem>
                    <MenuItem value={5}>RunningRainbow</MenuItem>
                    <MenuItem value={6}>RunningColor</MenuItem>
                </Select>
            </FormControl>
            <Typography className={props.classes.subtitle} gutterBottom>
                Color
            </Typography>
            <HuePicker
                color={props.color}
                onChange={props.onColorChange}
                onChangeComplete={props.onChange}
            />
            <div style={{margin: 16}}>
                <SliderPicker
                    color={props.color}
                    onChange={props.onColorChange}
                    onChangeComplete={props.onChange}
                />
            </div>
        </CardContent>
    </Card>;
}

LEDConfig.propTypes = {
    classes: PropTypes.any,
    onClick: PropTypes.func,
    onClick1: PropTypes.func,
    checked: PropTypes.bool,
    prop4: PropTypes.func,
    value: PropTypes.number,
    prop6: PropTypes.func,
    color: PropTypes.string,
    onChange: PropTypes.func
};
