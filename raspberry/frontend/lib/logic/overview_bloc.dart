import 'dart:async';
import 'dart:ui';

import 'package:bloc/bloc.dart';
import 'package:frontend/data/Repository.dart';
import 'package:meta/meta.dart';

@immutable
abstract class OverviewEvent {}

class InitEvent extends OverviewEvent {}

class SelectColorEvent extends OverviewEvent {
  final Color color;

  SelectColorEvent(this.color);
}

class SelectEffectEvent extends OverviewEvent {
  final int effectId;

  SelectEffectEvent(this.effectId);
}

class SelectPaletteEvent extends OverviewEvent {
  final int paletteId;

  SelectPaletteEvent(this.paletteId);
}

class OverviewState {
  final Color pickerColor;
  final List<String> effects;
  final List<String> palettes;

  OverviewState(this.pickerColor, this.effects, this.palettes);
}

class InitState extends OverviewState {
  InitState() : super(Color(0xff443a49), [], []);
}

class OverviewBloc extends Bloc<OverviewEvent, OverviewState> {
  final Repository repository;
  Color color = Color(0xff443a49);

  OverviewBloc({@required this.repository})
      : assert(repository != null),
        super(InitState());

  @override
  Stream<OverviewState> mapEventToState(
    OverviewEvent event,
  ) async* {
    switch (event.runtimeType) {
      case InitEvent:
        yield OverviewState(color, await repository.getNodeEffects(),
            await repository.getNodePalettes());
        break;
      case SelectColorEvent:
        color = (event as SelectColorEvent).color;
        repository.setNodeColor(color);
        yield OverviewState(color, await repository.getNodeEffects(),
            await repository.getNodePalettes());
        break;
      case SelectEffectEvent:
        repository.selectNodeEffect((event as SelectEffectEvent).effectId);
        break;
      case SelectPaletteEvent:
        repository.selectNodePalette((event as SelectPaletteEvent).paletteId);
        break;
    }
  }
}
