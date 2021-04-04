import 'dart:ui';

import 'package:frontend/data/api.dart';
import 'package:frontend/data/node_api.dart';
import 'package:frontend/data/node_model.dart';

class Repository {
  List<NodeModel> _nodes;

  Future<List<NodeModel>> getNodes() async {
    if (_nodes == null) {
      _nodes = (await fetchNodes()).nodes;
    }
    return _nodes;
  }

  Future<List<String>> getNodeEffects() async {
    var nodes = await getNodes();
    if (nodes.isNotEmpty) {
      return nodes[0].effects;
    } else {
      return [];
    }
  }

  Future<List<String>> getNodePalettes() async {
    var nodes = await getNodes();
    if (nodes.isNotEmpty) {
      return nodes[0].palettes;
    } else {
      return [];
    }
  }

  void selectNodeEffect(int effectId) async {
    (await getNodes()).where((node) => node.isSelected).forEach((node) {
      node.setEffect(effectId);
    });
  }

  void selectNodePalette(int paletteId) async {
    (await getNodes()).where((node) => node.isSelected).forEach((node) {
      node.setPalette(paletteId);
    });
  }

  void setNodeColor(Color color) async {
    (await getNodes()).where((node) => node.isSelected).forEach((node) {
      node.setColor(color);
    });
  }
}
