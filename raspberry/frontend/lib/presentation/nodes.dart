import 'package:flutter/material.dart';
import 'package:flutter_bloc/flutter_bloc.dart';
import 'package:frontend/logic/node_bloc.dart';
import 'package:frontend/extensions.dart';

class NodesList extends StatelessWidget {
  @override
  Widget build(BuildContext context) {
    return BlocBuilder<NodeBloc, NodeStates>(
        builder: (context, state) => Column(
            children:
                _buildNodeElement(BlocProvider.of<NodeBloc>(context), state)));
  }

  _buildNodeElement(NodeBloc bloc, NodeStates state) {
    List<Widget> widgets = [];
    widgets.add(Padding(
        padding: EdgeInsets.all(8),
        child: Text(
          "Nodes",
          style: TextStyle(fontSize: 30),
        )));
    widgets.add(_buildButton("Refresh", () {
      bloc.add(RefreshEvent());
    }));

    widgets.add(_buildNodeRowFromNode(state.allState, true, -1, bloc));

    widgets.addAll(state.nodeStates
        .mapIndexed((index, nodeState) =>
            _buildNodeRowFromNode(nodeState, false, index, bloc))
        );

    return widgets;
  }

  Widget _buildButton(String name, Function callback) => ElevatedButton(
      onPressed: callback,
      child: Container(
          width: 64,
          height: 64,
          child: Center(
            child: Text(name),
          )));

  Widget _buildNodeRowFromNode(
          NodeState state, bool all, int id, NodeBloc bloc) =>
      _buildNodeRow(
        state.name,
        state.brightness.toDouble(),
        (brightness) {
          bloc.add(BrightnessEvent(brightness, all, id));
        },
        state.on,
        (on) {
          bloc.add(OnOffEvent(on, all, id));
        },
        state.isSelected,
        (isSelected) {
          bloc.add(SelectEvent(isSelected, all, id));
        },
      );

  Widget _buildNodeRow(
          String name,
          double briValue,
          Function onBriChange,
          bool onValue,
          Function onOnChange,
          bool checked,
          Function onCheckedChange) =>
      Row(
        mainAxisAlignment: MainAxisAlignment.center,
        children: [
          Padding(
            padding: const EdgeInsets.all(8),
            child: Text(name),
          ),
          Slider(value: briValue, min: 0, max: 255, onChanged: onBriChange),
          Switch(value: onValue, onChanged: onOnChange),
          Checkbox(value: checked, onChanged: onCheckedChange)
        ],
      );
}
