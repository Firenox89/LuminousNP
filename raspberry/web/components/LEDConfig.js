import Typography from "@material-ui/core/Typography";
import FormControlLabel from "@material-ui/core/FormControlLabel";
import Checkbox from "@material-ui/core/Checkbox";
import FormControl from "@material-ui/core/FormControl";
import InputLabel from "@material-ui/core/InputLabel";
import Select from "@material-ui/core/Select";
import MenuItem from "@material-ui/core/MenuItem";
import {SliderPicker} from "react-color";
import * as PropTypes from "prop-types";
import React from "react";
import Paper from "@material-ui/core/Paper";
import {Button} from "@material-ui/core";

export function LEDConfig(props) {
    const [effectList, setEffectList] = React.useState();

    if (!effectList) {
        fetch("/getEffectList")
            .then(response => response.json())
            .then(data => {
                setEffectList(data)
            })
    }

    const buildEffectMenuItems = () => {
        if (effectList) {
            return effectList.map(value => {
                return (
                    <MenuItem key={value.ID} value={value.ID}>{value.Name}</MenuItem>
                )
            })
        }
    }

    const needsColor = () => {
        if (effectList) {
            const effect = effectList.find(value => value.ID === props.selectedEffect)
            if (effect) {
                return effect.NeedsColor
            }
        }
        return false
    }

    return <Paper className={props.classes.paper}>
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
            <Select
                id="demo-simple-select-outlined"
                value={props.selectedEffect}
                onChange={(props.onEffectChange)}
            >
                {buildEffectMenuItems()}
            </Select>
        </FormControl>
        <Button variant="contained" className={props.classes.button} onClick={props.onApply}>Apply</Button>
        {/*<HuePicker*/}
        {/*    color={props.color}*/}
        {/*    onChange={props.onColorChange}*/}
        {/*    onChangeComplete={props.onChange}*/}
        {/*/>*/}
        {needsColor() &&
        <div>
            <Typography className={props.classes.subtitle} gutterBottom>
                Color
            </Typography>
            <div style={{margin: 16}}>
                <SliderPicker
                    color={props.color}
                    onChange={props.onColorChange}
                    onChangeComplete={props.onChange}
                />
            </div>
        </div>
        }
    </Paper>;
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
