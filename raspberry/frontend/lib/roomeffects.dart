import 'package:flutter/material.dart';
import 'package:frontend/api_slider.dart';
import 'package:http/http.dart' as http;
import 'dart:html';
import 'dart:convert';

class RoomEffects extends StatefulWidget {
  @override
  _RoomEffectsPageState createState() => _RoomEffectsPageState();
}

class _RoomEffectsPageState extends State<RoomEffects> {
  List<String> effects = [];
  List<String> palettes = [];

  _RoomEffectsPageState() {
    _loadData();
  }

  @override
  Widget build(BuildContext context) => Row(children: [
    Padding(
        padding: EdgeInsets.all(16),
        child: SingleChildScrollView(
            child: Column(children: _buildEffectSelector()))),
    Padding(
        padding: EdgeInsets.all(16),
        child: SingleChildScrollView(
            child: Column(children: _buildPaletteSelector()))),
    Padding(
        padding: EdgeInsets.all(16),
        child: SingleChildScrollView(
            child: Column(children: _buildPaletteSelector()))),
  ]);

  _onPaletteSelect(int id) async {
    _sendRequest("setPalette", {'id': id.toString()});
  }

  List<Widget> _buildEffectSettings() {
    List<Widget> widgets = [];

    widgets.add(Padding(
        padding: EdgeInsets.all(8),
        child: Text(
          "Settings",
          style: TextStyle(fontSize: 30),
        )));

    widgets.add(Padding(
        padding: EdgeInsets.all(8),
        child: ApiSlider(
          title: "Brightness",
          urlApi: "/brightness",
          min: 0,
          max: 100
        )));
    return widgets;
  }

  List<Widget> _buildPaletteSelector() {
    if (palettes.isNotEmpty) {
      return _buildListSelector(palettes, _onPaletteSelect, "Palette");
    }
    return [];
  }

  _onEffectSelect(int id) async {
    _sendRequest("setEffect", {'id': id.toString()});
  }

  List<Widget> _buildEffectSelector() {
    if (effects.isNotEmpty) {
      return _buildListSelector(effects, _onEffectSelect, "Effect");
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

  Widget _buildButton(String name, Function callback) => ElevatedButton(
      onPressed: callback,
      child: Container(
          width: 64,
          height: 64,
          child: Center(
            child: Text(name),
          )));

  Future<void> _loadData() async {
    print("load data");

    effects = await _sendRequest("getEffectList");
    print("effects");
    palettes = await _sendRequest("getColorPaletteList");

    setState(() {});
  }

  _sendRequest(String endpoint, [Map<String, String> query]) async {
    String data;
    String _ref = window.location.href;

    try {
      data = await http.read(Uri.http(_ref, endpoint, query));
    } catch (error) {
      data =
      await http.read(Uri.http("localhost:1234", endpoint, query));
    }
    print(data);
    return jsonDecode(data).cast<String>();
  }
}
