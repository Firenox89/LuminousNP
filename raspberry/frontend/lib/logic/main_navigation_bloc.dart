import 'dart:async';

import 'package:bloc/bloc.dart';
import 'package:meta/meta.dart';

enum NavEvent {Overview, RoomEffects, Settings}

@immutable
abstract class MainNavigationState {}

class MainNavigationOverview extends MainNavigationState {}
class MainNavigationRoomEffects extends MainNavigationState {}
class MainNavigationSettings extends MainNavigationState {}

class MainNavigationBloc
    extends Bloc<NavEvent, MainNavigationState> {
  MainNavigationBloc() : super(MainNavigationOverview());

  @override
  Stream<MainNavigationState> mapEventToState(
    NavEvent event,
  ) async* {
    switch (event) {
      case NavEvent.Overview:
        yield MainNavigationOverview();
        break;
      case NavEvent.RoomEffects:
        yield MainNavigationRoomEffects();
        break;
      case NavEvent.Settings:
        yield MainNavigationSettings();
        break;
    }
  }
}
