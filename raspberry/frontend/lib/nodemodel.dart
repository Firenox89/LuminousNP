import 'dart:ui';

import 'package:http/http.dart' as http;

class NodeModel {
  List<ConnectedMCUs> connectedMCUs;

  NodeModel({this.connectedMCUs});

  NodeModel.fromJson(Map<String, dynamic> json) {
    if (json['ConnectedMCUs'] != null) {
      connectedMCUs = new List<ConnectedMCUs>();
      json['ConnectedMCUs'].forEach((v) {
        connectedMCUs.add(new ConnectedMCUs.fromJson(v));
      });
    }
  }
}

class ConnectedMCUs {
  String iP;
  String iD;
  int type;
  List<String> effects;
  List<String> palettes;
  int brightness;
  bool on;
  bool isSelected = true;

  ConnectedMCUs.fromJson(Map<String, dynamic> json) {
    iP = json['IP'];
    iD = json['ID'];
    type = json['Type'];
    effects = json['Effects'].cast<String>();
    palettes = json['Palettes'].cast<String>();
    brightness = json['Brightness'];
    on = json['On'];

    print("init bri " + brightness.toString());
  }

  Future<bool> toggleOnOff(bool on) async {
    await http
        .read(Uri.http(iP, "win&T=" + (on ? 1 : 0).toString() + "2&SN=0"));
    this.on = on;
    return on;
  }

  Future<int> setBrightness(int brightness) async {
    await http.read(Uri.http(iP, "win&A=" + brightness.toString() + "&SN=0"));
    this.brightness = brightness;
    return brightness;
  }

  setEffect(int effectId) async {
    await http.read(Uri.http(iP, "win&FX=" + effectId.toString() + "&SN=0"));
  }

  setPalette(int id) async {
    await http.read(Uri.http(iP, "win&FP=" + id.toString() + "&SN=0"));
  }

  setColor(Color color) async {
    await http.read(Uri.http(
        iP,
        "win&R=" +
            color.red.toString() +
            "&B=" +
            color.blue.toString() +
            "&G=" +
            color.green.toString() +
            "&SN=0"));
  }
}
