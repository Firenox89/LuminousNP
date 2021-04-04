import 'package:flutter/cupertino.dart';
import 'package:flutter/material.dart';
import 'package:flutter/rendering.dart';
import 'package:flex_color_picker/flex_color_picker.dart';
import 'package:flutter_bloc/flutter_bloc.dart';
import 'package:frontend/logic/overview_bloc.dart';
import 'package:frontend/presentation/nodes.dart';

class OverviewPage extends StatelessWidget {
  @override
  Widget build(BuildContext context) =>
      BlocBuilder<OverviewBloc, OverviewState>(
          builder: (context, state) => Row(children: [
                Container(
                    decoration: BoxDecoration(
                        border: Border(right: BorderSide(color: Colors.white))),
                    child: Padding(
                        padding: EdgeInsets.all(16), child: NodesList())),
                Container(
                    width: 400,
                    height: 400,
                    child: ColorPicker(
                      color: state.pickerColor,
                      onColorChanged: (color) {
                        BlocProvider.of<OverviewBloc>(context)
                            .add(SelectColorEvent(color));
                      },
                    )),
                Padding(
                    padding: EdgeInsets.all(16),
                    child: SingleChildScrollView(
                        child: Column(
                            children: _buildEffectSelector(
                                BlocProvider.of<OverviewBloc>(context),
                                state)))),
                Padding(
                    padding: EdgeInsets.all(16),
                    child: SingleChildScrollView(
                        child: Column(
                            children: _buildPaletteSelector(
                                BlocProvider.of<OverviewBloc>(context),
                                state)))),
              ]));

  List<Widget> _buildPaletteSelector(OverviewBloc bloc, OverviewState state) {
    return _buildListSelector(state.palettes, (index) {
      bloc.add(SelectPaletteEvent(index));
    }, "Color");
  }

  List<Widget> _buildEffectSelector(OverviewBloc bloc, OverviewState state) {
    return _buildListSelector(state.effects, (index) {
      bloc.add(SelectEffectEvent(index));
    }, "Effects");
  }

  List<Widget> _buildListSelector(
      List<String> list, Function callback, String title) {
    List<Widget> widgets = [];
    widgets.add(Padding(
        padding: EdgeInsets.all(8),
        child: Text(
          title,
          style: TextStyle(fontSize: 30),
        )));
    Map<int, Widget> buttonMap = list.asMap().map((key, value) => MapEntry(
        key,
        Padding(
            padding: EdgeInsets.all(8),
            child: ElevatedButton(
                onPressed: () {
                  callback(key);
                },
                child: Text(value)))));
    widgets.addAll(buttonMap.values);
    return widgets;
  }
}
