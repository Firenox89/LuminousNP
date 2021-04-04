class NodesModel {
  List<NodeModel> nodes = [];

  NodesModel({this.nodes});

  NodesModel.fromJson(List<dynamic> json) {
    json.forEach((v) {
      nodes.add(NodeModel.fromJson(v));
    });
  }
}

class NodeModel {
  String iP;
  String iD;
  int type;
  List<String> effects;
  List<String> palettes;
  int brightness;
  bool on;
  bool isSelected = true;

  NodeModel.fromJson(Map<String, dynamic> json) {
    iP = json['IP'];
    iD = json['ID'];
    type = json['Type'];
    effects = json['Effects'].cast<String>();
    palettes = json['Palettes'].cast<String>();
    brightness = json['Brightness'];
    on = json['On'];
  }
}
