import 'package:flutter/material.dart';
import 'package:frontend/nodemodel.dart';
import 'package:frontend/wled-api.dart';
import 'package:http/http.dart' as http;
import 'dart:html';
import 'dart:convert';

void main() {
  runApp(MyApp());
}

class MyApp extends StatelessWidget {
  @override
  Widget build(BuildContext context) {
    return MaterialApp(
      title: 'Flutter Demo',
      theme: ThemeData.dark(),
      home: MyHomePage(title: 'Flutter Demo Home Page'),
    );
  }
}

class MyHomePage extends StatefulWidget {
  MyHomePage({Key key, this.title}) : super(key: key);
  final String title;

  @override
  _MyHomePageState createState() => _MyHomePageState();
}

class _MyHomePageState extends State<MyHomePage> {
  List<ConnectedMCUs> nodes = [];
  List<String> effects = [];
  List<String> palettes = [];

  _MyHomePageState() {
    _loadNodes();
  }

  Future<void> _loadNodes() async {
    String _ref = window.location.href;
    String data;
    try {
      data = await http.read(Uri.http(_ref, "getConnectedNodeMCUs"));
    } catch (error) {
      data =
          await http.read(Uri.http("localhost:1234", "getConnectedNodeMCUs"));
    }

    var json = NodeModel.fromJson(jsonDecode(data));
    nodes = json.connectedMCUs;
    effects = nodes.first.effects;
    palettes = nodes.first.palettes;

    setState((){});
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(
        title: Text(widget.title),
      ),
      body: Center(
        child: Column(
          mainAxisAlignment: MainAxisAlignment.center,
          children: <Widget>[
            ElevatedButton(onPressed: _loadNodes, child: Text("Refresh")),
            Container(
                padding: EdgeInsets.only(top: 25),
                height: 500,
                width: 250,
                child: ListView.builder(
                    itemCount: nodes.length,
                    itemBuilder: (BuildContext context, int index) {
                      return buildNodeRow(nodes[index]);
                    })),
          ],
        ),
      ),
    );
  }

  Widget buildNodeRow(ConnectedMCUs node) {
    return Row(mainAxisAlignment: MainAxisAlignment.center, children: [
      Padding(
        padding: const EdgeInsets.all(8.0),
        child: Text(node.iD),
      ),
      ElevatedButton(onPressed: () {toggleOnOff(node);}, child: Text("On"))
    ], );
  }
}
