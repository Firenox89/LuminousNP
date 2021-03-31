import 'package:flutter/cupertino.dart';
import 'package:flutter/material.dart';
import 'package:flutter/rendering.dart';
import 'package:frontend/api.dart';
import 'dart:convert';
import 'package:frontend/nodemodel.dart';
import 'package:flex_color_picker/flex_color_picker.dart';

class OverviewPage extends StatefulWidget {
  @override
  _OverviewPageState createState() => _OverviewPageState();
}

class _OverviewPageState extends State<OverviewPage> {
  List<ConnectedMCUs> nodes = [];
  List<String> effects = [];
  List<String> palettes = [];

  _OverviewPageState() {
    _loadNodes();
  }

  Color pickerColor = Color(0xff443a49);

  @override
  Widget build(BuildContext context) => Row(children: [
        Container(
            decoration: BoxDecoration(
                border: Border(right: BorderSide(color: Colors.white))),
            child: Padding(
                padding: EdgeInsets.all(16),
                child: Column(children: _buildNodeListWidget()))),
        Container(
            width: 400,
            height: 400,
            child: ColorPicker(
              color: pickerColor,
              onColorChanged: (color) {
                changeColor(color);
              },
            )),
        Padding(
            padding: EdgeInsets.all(16),
            child: SingleChildScrollView(
                child: Column(children: _buildEffectSelector()))),
        Padding(
            padding: EdgeInsets.all(16),
            child: SingleChildScrollView(
                child: Column(children: _buildPaletteSelector()))),
      ]);

  void changeColor(Color color) {
    nodes.where((element) => element.isSelected).forEach((element) {
      element.setColor(color);
    });
    setState(() => pickerColor = color);
  }

  _onPaletteSelect(int id) {
    nodes.where((element) => element.isSelected).forEach((element) {
      element.setPalette(id);
    });
  }

  List<Widget> _buildPaletteSelector() {
    if (nodes.isNotEmpty) {
      return _buildListSelector(nodes[0].palettes, _onPaletteSelect, "Color");
    }
    return [];
  }

  _onEffectSelect(int id) {
    nodes.where((element) => element.isSelected).forEach((element) {
      element.setEffect(id);
    });
  }

  List<Widget> _buildEffectSelector() {
    if (nodes.isNotEmpty) {
      return _buildListSelector(nodes[0].effects, _onEffectSelect, "Effects");
    }
    return [];
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

  List<Widget> _buildNodeListWidget() {
    List<Widget> widgets = [];
    widgets.add(Padding(
        padding: EdgeInsets.all(8),
        child: Text(
          "Nodes",
          style: TextStyle(fontSize: 30),
        )));
    widgets.add(_buildButton("Refresh", () {
      _loadNodes();
    }));
    if (nodes.isNotEmpty) {
      widgets.add(_buildNodeRow(
          "All",
          nodes[0].brightness.toDouble(),
          _setBrightnessOnAllNodes,
          nodes[0].on,
          _setOnOnAllNodes,
          nodes[0].isSelected, (value) {
        nodes.forEach((element) {
          element.isSelected = value;
        });
        setState(() {});
      }));
    }
    widgets.addAll(nodes.map((e) => _buildNodeRowFromNode(e)));

    return widgets;
  }

  _setBrightnessOnAllNodes(double brightness) async {
    await Future.wait(nodes.map((element) {
      return element.setBrightness(brightness.toInt());
    }));
    setState(() {});
  }

  _setOnOnAllNodes(bool on) async {
    await Future.wait(nodes.map((element) {
      return element.toggleOnOff(on);
    }));
    setState(() {});
  }

  Widget _buildButton(String name, Function callback) => ElevatedButton(
      onPressed: callback,
      child: Container(
          width: 64,
          height: 64,
          child: Center(
            child: Text(name),
          )));

  Future<void> _loadNodes() async {
    print("load nodes");

    var data = await request("getConnectedNodeMCUs");
    var json = NodeModel.fromJson(jsonDecode(data));
    nodes = json.connectedMCUs;
    effects = nodes.first.effects;
    palettes = nodes.first.palettes;

    setState(() {});
  }

  Widget _buildNodeRowFromNode(ConnectedMCUs node) => _buildNodeRow(
      node.iD,
      node.brightness.toDouble(),
      (brightness) {
        node.setBrightness(brightness.toInt()).then((value) => setState(() {}));
      },
      node.on,
      (on) {
        print("switch " + on.toString());
        node.toggleOnOff(on).then((value) => setState(() {}));
      },
      node.isSelected,
      (value) {
        node.isSelected = value;
        setState(() {});
      });

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
