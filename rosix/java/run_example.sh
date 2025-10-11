#!/bin/bash

echo "=== 编译Java示例 ==="
javac -cp ".:lib/*" \
  -d target/classes \
  src/main/java/com/uros/rosix/example/RealWorldExample.java \
  src/main/java/com/uros/rosix/core/*.java

echo ""
echo "=== 运行示例 ==="
java -cp "target/classes:lib/*" \
  com.uros.rosix.example.RealWorldExample
