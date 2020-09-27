import Typography from "@material-ui/core/Typography";
import {HuePicker} from "react-color";
import * as PropTypes from "prop-types";
import React from "react";
import Paper from "@material-ui/core/Paper";
import {Button} from "@material-ui/core";
import Grid from "@material-ui/core/Grid";
import Slider from "@material-ui/core/Slider";
import {ColorPaletteButton} from "./ColorPaletteButton";

export function LEDConfig(props) {
    const [effectList, setEffectList] = React.useState();
    const [colorPaletteList, setColorPaletteList] = React.useState();

    if (!effectList) {
        fetch("/getEffectList")
            .then(response => response.json())
            .then(data => {
                setEffectList(data)
            })
    }
    if (!colorPaletteList) {
        fetch("/getColorPaletteList")
            .then(response => response.json())
            .then(data => {
                setColorPaletteList(data)
            })
    }

    const buildEffectTiles = () => {
        if (effectList) {
            return effectList.map(value => {
                return (
                    <Button key={value.ID} variant="contained" className={props.classes.button}
                            onClick={() => props.onEffectChange(value.ID)}>{value.Name}</Button>
                )
            })
        }
    }

    const buildColorPaletteTiles = () => {
        if (colorPaletteList) {
            return colorPaletteList.map(value => {
                return (
                    <ColorPaletteButton key={value.ID} classes={props.classes}
                                        onClick={() => props.onColorPaletteChange(value.ID)}
                                        colors={value.Colors}
                    />
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

    const needsColorPalette = () => {
        if (effectList) {
            const effect = effectList.find(value => value.ID === props.selectedEffect)
            if (effect) {
                return effect.NeedsColorPalette
            }
        }
        return false
    }

    return <Paper className={props.classes.paper}>
        <Typography className={props.classes.title} gutterBottom>
            Preset
        </Typography>
        <Grid container spacing={1}>
            <Grid item xs={12} sm={12}>
                <Button variant="contained" className={props.classes.button} onClick={props.onOff}>Off</Button>
                {buildEffectTiles()}
            </Grid>
            {needsColorPalette() &&
            <Grid item xs={12} sm={12}>
                {buildColorPaletteTiles()}
            </Grid>
            }
        </Grid>
        {needsColor() &&
        <div>
            <Typography className={props.classes.subtitle} gutterBottom>
                Color
            </Typography>
            <div style={{margin: 16}}>
                <HuePicker
                    color={props.color}
                    onChange={props.onColorChange}
                    onChangeComplete={props.onChange}
                />
            </div>
        </div>
        }
        <Typography className={props.classes.title} gutterBottom>
            Brightness
        </Typography>
        <div style={{margin: 16}}>
            <Slider value={props.brightness} onChange={(event, newValue) => props.onChangeBrightness(newValue)}/>
        </div>
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
    onChange: PropTypes.func,
    brightness: PropTypes.number,
    onChangeBrightness: PropTypes.func
};
