extension ExtendedIterable<E> on Iterable<E> {
  /// Like Iterable<T>.map but callback have index as second argument
  Iterable<T> mapIndexed<T>(T f(int i, E e)) {
    var i = 0;
    return this.map((e) => f(i++, e)).toList();
  }

  void forEachIndexed(void f(E e, int i)) {
    var i = 0;
    this.forEach((e) => f(e, i++));
  }
}