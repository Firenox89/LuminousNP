import {Button} from "@material-ui/core";
import React from "react";
import CircularBar from 'react-multicolor-circular-progress-bar';

export function ColorPaletteButton(props) {
    let angles = []
    let inc = 360 / props.colors.length
    for (let i = 1; i <= props.colors.length-1; i++) {
        angles.push((i * inc)-1)
    }
    return (
        <Button variant="contained" className={props.classes.button} onClick={props.onClick}>
            <CircularBar
                scale={0.3}
                angleTransition={angles}
                colors={props.colors}
                stroke={{color: '#eee', width: 102}}
                percent={{value:100, showValue:false}}
            />
        </Button>
    )
}
