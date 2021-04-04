import 'dart:convert';

import 'package:frontend/data/node_model.dart';
import 'package:http/http.dart' as http;
import 'dart:html';

Future<String> fetch(String endpoint) async {
  String _ref = window.location.href;
  String data;
  try {
    data = await http.read(Uri.http(_ref, endpoint));
  } catch (error) {
    data = await http.read(Uri.http("localhost:1234", endpoint));
  }
  return data;
}

Future<NodesModel> fetchNodes() async {
  var data = await fetch("getConnectedNodeMCUs");
  return NodesModel.fromJson(jsonDecode(data));
}