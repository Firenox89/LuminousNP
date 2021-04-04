import 'package:flutter/material.dart';
import 'package:frontend/data/node_model.dart';
import 'package:http/http.dart' as http;

extension ApiExtendedNodeModel on NodeModel {
  Future<bool> setOn(bool on) async {
    await http
        .read(Uri.http(this.iP, "win&T=" + (on ? 1 : 0).toString() + "2&SN=0"));
    this.on = on;
    return on;
  }

  Future<int> setBrightness(int brightness) async {
    await http
        .read(Uri.http(this.iP, "win&A=" + brightness.toString() + "&SN=0"));
    this.brightness = brightness;
    return brightness;
  }

  setEffect(int effectId) async {
    await http
        .read(Uri.http(this.iP, "win&FX=" + effectId.toString() + "&SN=0"));
  }

  setPalette(int id) async {
    await http.read(Uri.http(this.iP, "win&FP=" + id.toString() + "&SN=0"));
  }

  setColor(Color color) async {
    await http.read(Uri.http(
        this.iP,
        "win&R=" +
            color.red.toString() +
            "&B=" +
            color.blue.toString() +
            "&G=" +
            color.green.toString() +
            "&SN=0"));
  }
}
