<?php

$iter = new DirectoryIterator(__DIR__);

foreach ($iter as $obj) {
  if ($obj->isDot() || $obj->isDir()) {
    continue;
  }

  $name = $obj->getFilename();

  if (substr($name, -4) !== 'json') {
    continue;
  }

  $data = file_get_contents($name);
  $lines = count(explode("\n", $data));

  if ($lines <= 1) {
    $decoded = json_decode($data, true);
    $encoded = json_encode($decoded, JSON_PRETTY_PRINT);
    $encoded = preg_replace('/^(  +?)\\1(?=[^ ])/m', '$1', $encoded);

    file_put_contents($name, $encoded);
    echo "UPDATED: $name\n";
  }
}

echo "DONE?\n";
