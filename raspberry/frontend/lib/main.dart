import 'package:flutter/material.dart';
import 'package:flutter_bloc/flutter_bloc.dart';
import 'package:frontend/data/Repository.dart';
import 'package:frontend/logic/main_navigation_bloc.dart';
import 'package:frontend/logic/node_bloc.dart';
import 'package:frontend/logic/overview_bloc.dart';
import 'package:frontend/presentation/overview.dart';
import 'package:frontend/presentation/room_effects.dart';
import 'package:frontend/presentation/settings.dart';

void main() {
  var repo = Repository();
  runApp(
    MultiBlocProvider(
        providers: [
          BlocProvider(create: (context) => MainNavigationBloc()),
          BlocProvider(create: (context) => OverviewBloc(repository: repo)..add(InitEvent())),
          BlocProvider(create: (context) => NodeBloc(repository: repo)..add(RefreshEvent())),
        ],
        child: MaterialApp(
          theme: ThemeData.dark(),
          home: Scaffold(body: App()),
        )),
  );
}

class App extends StatelessWidget {
  @override
  Widget build(BuildContext context) {
    return BlocBuilder<MainNavigationBloc, MainNavigationState>(
        builder: (context, state) => Row(
              children: [
                Container(
                    width: 90,
                    decoration: BoxDecoration(
                        border: Border(right: BorderSide(color: Colors.white))),
                    child: Column(
                      mainAxisAlignment: MainAxisAlignment.spaceEvenly,
                      children: _buildNavigation(
                          BlocProvider.of<MainNavigationBloc>(context)),
                    )),
                Container(
                  child: _buildNavContent(state),
                ),
              ],
            ));
  }

  List<Widget> _buildNavigation(MainNavigationBloc bloc) {
    return [
      _buildNavButton("Overview", () {
        bloc.add(NavEvent.Overview);
      }),
      _buildNavButton("Room Effects", () {
        bloc.add(NavEvent.RoomEffects);
      }),
      _buildNavButton("Settings", () {
        bloc.add(NavEvent.Settings);
      }),
    ]
        .map((e) => Padding(
              padding: const EdgeInsets.all(8.0),
              child: e,
            ))
        .toList();
  }

  Widget _buildNavButton(String name, Function onPressed) => ElevatedButton(
      onPressed: onPressed,
      child: Container(
          width: 64,
          height: 64,
          child: Center(
            child: Text(name),
          )));

  Widget _buildNavContent(MainNavigationState state) {
    switch (state.runtimeType) {
      case MainNavigationOverview:
        return OverviewPage();
        break;
      case MainNavigationRoomEffects:
        return RoomEffects();
        break;
      case MainNavigationSettings:
        return Settings();
        break;
    }
    throw Exception("Boom");
  }
}
