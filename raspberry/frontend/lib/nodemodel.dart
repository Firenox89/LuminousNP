class NodeModel {
  List<ConnectedMCUs> connectedMCUs;

  NodeModel({this.connectedMCUs});

  NodeModel.fromJson(Map<String, dynamic> json) {
    if (json['ConnectedMCUs'] != null) {
      connectedMCUs = <ConnectedMCUs>[];
      json['ConnectedMCUs'].forEach((v) {
        connectedMCUs.add(new ConnectedMCUs.fromJson(v));
      });
    }
  }
}

class ConnectedMCUs {
  String iP;
  String iD;
  List<String> effects;
  List<String> palettes;

  ConnectedMCUs({this.iP, this.iD});

  ConnectedMCUs.fromJson(Map<String, dynamic> json) {
    iP = json['IP'];
    iD = json['ID'];
    effects = json['effects'];
    palettes = json['palettes'];
  }
}