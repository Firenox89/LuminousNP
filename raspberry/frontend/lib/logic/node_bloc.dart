import 'dart:async';

import 'package:bloc/bloc.dart';
import 'package:frontend/data/Repository.dart';
import 'package:frontend/data/node_model.dart';
import 'package:frontend/data/node_api.dart';
import 'package:meta/meta.dart';
import 'package:frontend/extensions.dart';

@immutable
abstract class NodeEvent {}

class RefreshEvent extends NodeEvent {}

class SelectEvent extends NodeEvent {
  final bool all;
  final int index;
  final bool isSelected;

  SelectEvent(this.isSelected, this.all, this.index);
}

class OnOffEvent extends NodeEvent {
  final bool all;
  final int index;
  final bool on;

  OnOffEvent(this.on, this.all, this.index);
}

class BrightnessEvent extends NodeEvent {
  final bool all;
  final int index;
  final double brightness;

  BrightnessEvent(this.brightness, this.all, this.index);
}

class NodeState {
  final String name;
  final int brightness;
  final bool on;
  final bool isSelected;

  NodeState(this.name, this.brightness, this.on, this.isSelected);
}

class NodeStates {
  final NodeState allState;
  final List<NodeState> nodeStates;

  NodeStates(this.allState, this.nodeStates);
}

class InitNodeState extends NodeStates {
  InitNodeState() : super(NodeState("All", 0, false, false), []);
}

class NodeBloc extends Bloc<NodeEvent, NodeStates> {
  final Repository repository;

  NodeBloc({@required this.repository})
      : assert(repository != null),
        super(InitNodeState()) {}

  @override
  void onChange(Change<NodeStates> change) {
    print(change);
    super.onChange(change);
  }

  @override
  Stream<NodeStates> mapEventToState(NodeEvent event) async* {
    switch (event.runtimeType) {
      case RefreshEvent:
        yield await _handleRefreshEvent();
        break;
      case SelectEvent:
        yield await _handleSelectEvent(event);
        break;
      case OnOffEvent:
        yield await _handleOnOffEvent(event);
        break;
      case BrightnessEvent:
        yield await _handleBrightnessEvent(event);
        break;
    }
  }

  NodeStates _buildNodeStates(List<NodeModel> nodes) {
    if (nodes.isEmpty) {
      return NodeStates(NodeState("All", 0, false, false), []);
    } else {
      var allOn = nodes.any((node) => node.on);
      var allSelected = nodes.any((node) => node.isSelected);
      var allBrightness = nodes[0].brightness;
      return NodeStates(
          NodeState("All", allBrightness, allOn, allSelected),
          nodes.mapIndexed((index, node) => NodeState(
              "${index + 1}", node.brightness, node.on, node.isSelected)));
    }
  }

  _handleRefreshEvent() async {
    var nodes = await this.repository.getNodes();
    return _buildNodeStates(nodes);
  }

  _handleSelectEvent(SelectEvent event) async {
    var nodes = await this.repository.getNodes();
    if (event.all) {
      nodes.forEach((node) {
        node.isSelected = event.isSelected;
      });
    } else {
      nodes[event.index].isSelected = event.isSelected;
    }
    return _buildNodeStates(nodes);
  }

  _handleOnOffEvent(OnOffEvent event) async {
    var nodes = await this.repository.getNodes();
    if (event.all) {
      await Future.wait(nodes.map((node) => node.setOn(event.on)));
    } else {
      nodes[event.index].setOn(event.on);
    }
    return _buildNodeStates(nodes);
  }

  _handleBrightnessEvent(BrightnessEvent event) async {
    var nodes = await this.repository.getNodes();
    if (event.all) {
      await Future.wait(
          nodes.map((node) => node.setBrightness(event.brightness.toInt())));
    } else {
      nodes[event.index].setBrightness(event.brightness.toInt());
    }
    return _buildNodeStates(nodes);
  }
}
