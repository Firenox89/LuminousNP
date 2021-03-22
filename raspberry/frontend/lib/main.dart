import 'package:flutter/material.dart';
import 'package:frontend/overview.dart';
import 'package:frontend/roomeffects.dart';
import 'package:frontend/settings.dart';

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

enum Navigation { Overview, RoomEffects, Settings }

class _MyHomePageState extends State<MyHomePage> {
  Navigation currentNavigationItem = Navigation.Overview;

  @override
  Widget build(BuildContext context) {
    double deviceWidth = MediaQuery.of(context).size.width;
    double deviceHeight = MediaQuery.of(context).size.height;

    return Scaffold(
        body: Row(
          children: [
            Container(
                width: 90,
                height: deviceHeight,
                decoration: BoxDecoration(
                    border: Border(right: BorderSide(color: Colors.white))),
                child: Column(
                  mainAxisAlignment: MainAxisAlignment.spaceEvenly,
                  children: _buildNavigation(),
                )),
            Container(
              width: deviceWidth - 90,
              height: deviceHeight,
              child: _buildNavContent(),
            ),
          ],
        ));
  }

  _navigateTo(Navigation nav) {
    setState(() {
      currentNavigationItem = nav;
    });
  }

  Widget _buildNavButton(String name, Navigation navItem) => ElevatedButton(
      onPressed: () {
        _navigateTo(navItem);
      },
      child: Container(
          width: 64,
          height: 64,
          child: Center(
            child: Text(name),
          )));

  List<Widget> _buildNavigation() {
    return [
      _buildNavButton("Overview", Navigation.Overview),
      _buildNavButton("Room Effects", Navigation.RoomEffects),
      _buildNavButton("Settings", Navigation.Settings),
    ]
        .map((e) => Padding(
              padding: const EdgeInsets.all(8.0),
              child: e,
            ))
        .toList();
  }

  Widget _buildNavContent() {
    switch (currentNavigationItem) {
      case Navigation.Overview:
        return OverviewPage();
        break;
      case Navigation.RoomEffects:
        return RoomEffects();
        break;
      case Navigation.Settings:
        return Settings();
        break;
    }
    throw Exception("Boom");
  }
}
