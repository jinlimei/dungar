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

  if (strpos($name, '%!d') === false) {
    continue;
  }

  $new = str_replace(
    ['%!d(string=', ')'],
    '',
    $name
  );

  echo "NEED TO FIX $name (NEW $new )\n";
  rename($name, $new);
}

echo "DONE?\n";
