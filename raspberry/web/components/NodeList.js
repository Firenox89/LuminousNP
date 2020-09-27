import Typography from "@material-ui/core/Typography";
import Checkbox from "@material-ui/core/Checkbox";
import * as PropTypes from "prop-types";
import React from "react";
import Paper from "@material-ui/core/Paper";
import List from "@material-ui/core/List";
import ListItem from "@material-ui/core/ListItem";
import ListItemIcon from "@material-ui/core/ListItemIcon";
import ListItemText from "@material-ui/core/ListItemText";

export function NodeList(props) {
    const [allSelected, setAllSelected] = React.useState(false);

    const handleSelectAllClick = () => {
        if (allSelected) {
            setAllSelected(false)
            props.setSelected([]);
        } else {
            setAllSelected(true)
            const newSelecteds = props.connectedDevices.map((n) => n.ID);
            props.setSelected(newSelecteds);
        }
    };

    const handleClick = (name) => {
        const selectedIndex = props.selected.indexOf(name);
        let newSelected = [];

        if (selectedIndex === -1) {
            newSelected = newSelected.concat(props.selected, name);
        } else if (selectedIndex === 0) {
            newSelected = newSelected.concat(props.selected.slice(1));
        } else if (selectedIndex === props.selected.length - 1) {
            newSelected = newSelected.concat(props.selected.slice(0, -1));
        } else if (selectedIndex > 0) {
            newSelected = newSelected.concat(
                props.selected.slice(0, selectedIndex),
                props.selected.slice(selectedIndex + 1),
            );
        }

        props.setSelected(newSelected);
    };

    const handleSegmentClick = (name, index) => {

    };

    return <Paper className={props.classes.paper}>
        <Typography className={props.classes.title} gutterBottom>
            Nodes
        </Typography>
        <List>
            <ListItem key={"all"} role={undefined} dense button onClick={handleSelectAllClick}>
                <ListItemIcon>
                    <Checkbox
                        edge="start"
                        checked={allSelected}
                        tabIndex={-1}
                        onChange={handleSelectAllClick}
                    />
                </ListItemIcon>
                <ListItemText id={"allid"} primary={"All"}/>
            </ListItem>
            {props.connectedDevices.map((node) => {
                    return (
                        <ListItem key={node.ID} role={undefined} dense button onClick={() => handleClick(node.ID)}>
                            <ListItemIcon>
                                <Checkbox
                                    edge="start"
                                    checked={props.selected.indexOf(node.ID) !== -1}
                                    tabIndex={-1}
                                    onChange={() => handleClick(node.ID)}
                                    disabled={!node.IsConnected}
                                />
                            </ListItemIcon>
                            <ListItemText id={node.ID} primary={node.ID}/>
                            {/*<div>*/}
                            {/*    {node.Segments.map((segmentCount, index) => {*/}
                            {/*            return (<Checkbox*/}
                            {/*                key={index}*/}
                            {/*                edge="start"*/}
                            {/*                checked={props.selected.indexOf(node.ID) !== -1}*/}
                            {/*                tabIndex={-1}*/}
                            {/*                onChange={() => handleSegmentClick(node.ID, index)}*/}
                            {/*                disabled={!node.IsConnected}*/}
                            {/*            />)*/}
                            {/*        }*/}
                            {/*    )}*/}
                            {/*</div>*/}
                        </ListItem>
                    )
                }
            )}
        </List>
    </Paper>;
}

NodeList.propTypes = {
    classes: PropTypes.any,
    selected: PropTypes.array,
    onSelect: PropTypes.func,
    connectedDevices: PropTypes.array,
};
