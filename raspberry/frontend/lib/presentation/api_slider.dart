import 'package:flutter/material.dart';
import 'package:flutter_bloc/flutter_bloc.dart';

import 'package:http/http.dart' as http;

import '../data/api.dart';

class ApiSlider extends StatefulWidget {
  final String title;
  final String urlApi;
  final double min;
  final double max;

  ApiSlider({this.title, this.urlApi, this.min, this.max});

  @override
  _ApiSliderState createState() => _ApiSliderState();
}

class _ApiSliderState extends State<ApiSlider> {
  var value = 0.0;

  _ApiSliderState() {
    getValue();
  }

  getValue() async {
    value = (await fetch("brightness")) as double;
  }

  setValue(double newValue) async {
    value = (await fetch("brightness")) as double;
  }

  @override
  Widget build(BuildContext context) {
    Column(children: [
      Text(widget.title),
      Slider(value: value, min: widget.min, max: widget.max, onChanged: setValue),
    ]);
  }
}
