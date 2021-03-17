
import 'package:http/http.dart' as http;
import 'nodemodel.dart';

toggleOnOff(ConnectedMCUs node) async {
  print("toggle " + node.iD);
  await http.read(Uri.http(node.iP, "win&T=2&SN=0"));
}
