export interface TerritoryMapData {
  id: string;
  path: string;
  labelX: number;
  labelY: number;
}

export const territoryPaths: TerritoryMapData[] = [
  // === NORTH AMERICA === (x 25-360, y 25-330)
  // Continent shape: wide top (Alaska through Greenland), narrows down to Central America
  // Shared edges ensure puzzle-piece tiling with 1-2px gaps between territories

  {
    id: 'alaska',
    // Top-left of NA, bulky shape. Shares south-east edge with northwest_territory
    path: 'M 27 55 L 90 27 L 130 40 L 135 70 L 125 105 L 100 120 L 60 125 L 30 110 L 27 55 Z',
    labelX: 78,
    labelY: 78,
  },
  {
    id: 'northwest_territory',
    // Wide territory across top of NA. Shares west edge with alaska, south edges with alberta/ontario, east edge transitions toward greenland
    path: 'M 137 40 L 200 28 L 270 30 L 290 50 L 270 80 L 230 100 L 175 110 L 127 107 L 137 70 L 137 40 Z',
    labelX: 205,
    labelY: 68,
  },
  {
    id: 'greenland',
    // Island territory - separated by water gap from NW territory and quebec. Distinct island shape
    path: 'M 300 25 L 345 18 L 380 25 L 395 50 L 385 80 L 355 95 L 320 90 L 295 70 L 290 45 L 300 25 Z',
    labelX: 342,
    labelY: 55,
  },
  {
    id: 'alberta',
    // Below NW territory on the west side. Shares north edge with NW territory, east edge with ontario, south edge with western_us
    path: 'M 127 109 L 175 112 L 178 140 L 178 175 L 145 185 L 105 180 L 90 155 L 95 125 L 127 109 Z',
    labelX: 138,
    labelY: 148,
  },
  {
    id: 'ontario',
    // Center of NA. Shares west with alberta, north with NW territory, east with quebec, south with eastern_us/western_us
    path: 'M 177 112 L 232 102 L 272 82 L 292 95 L 288 130 L 270 160 L 230 175 L 180 177 L 180 140 L 177 112 Z',
    labelX: 232,
    labelY: 138,
  },
  {
    id: 'quebec',
    // East side of NA below greenland. Shares west with ontario, south with eastern_us
    path: 'M 272 84 L 294 70 L 325 78 L 350 100 L 345 135 L 325 160 L 295 170 L 272 162 L 290 130 L 274 84 Z',
    labelX: 312,
    labelY: 120,
  },
  {
    id: 'western_us',
    // Below alberta, west of eastern_us. Shares north with alberta, east with eastern_us, south with central_america
    path: 'M 95 182 L 147 187 L 180 179 L 185 210 L 175 245 L 150 265 L 115 268 L 80 250 L 70 220 L 80 195 L 95 182 Z',
    labelX: 130,
    labelY: 225,
  },
  {
    id: 'eastern_us',
    // East of western_us, below ontario/quebec. Shares west with western_us, north with ontario, south-east coast
    path: 'M 182 179 L 232 177 L 273 164 L 297 172 L 300 200 L 285 235 L 255 260 L 210 270 L 177 247 L 187 210 L 182 179 Z',
    labelX: 240,
    labelY: 218,
  },
  {
    id: 'central_america',
    // Narrow bottom of NA connecting to South America. Shares north with western_us/eastern_us
    path: 'M 117 270 L 152 267 L 178 249 L 212 272 L 210 295 L 190 318 L 165 328 L 140 320 L 120 300 L 110 280 L 117 270 Z',
    labelX: 162,
    labelY: 298,
  },

  // === SOUTH AMERICA === (x 130-320, y 330-570)
  // Classic shape: wide at top (Brazil), narrows to Argentina

  {
    id: 'venezuela',
    // Top of South America, wide. Shares south with peru and brazil
    path: 'M 145 332 L 195 332 L 250 338 L 290 350 L 280 375 L 245 385 L 205 388 L 165 380 L 140 360 L 145 332 Z',
    labelX: 215,
    labelY: 358,
  },
  {
    id: 'peru',
    // West side of SA. Shares north with venezuela, east with brazil, south with argentina
    path: 'M 140 362 L 167 382 L 207 390 L 210 420 L 200 450 L 175 468 L 148 460 L 132 430 L 130 395 L 140 362 Z',
    labelX: 170,
    labelY: 420,
  },
  {
    id: 'brazil',
    // Largest territory in SA - wide bulge. Shares west with venezuela/peru, south with argentina
    path: 'M 209 390 L 247 387 L 282 377 L 315 390 L 320 425 L 305 460 L 275 478 L 240 482 L 210 470 L 202 452 L 212 420 L 209 390 Z',
    labelX: 265,
    labelY: 430,
  },
  {
    id: 'argentina',
    // Bottom of SA, narrows to a point. Shares north with peru and brazil
    path: 'M 150 462 L 177 470 L 212 472 L 242 484 L 277 480 L 265 510 L 245 540 L 220 560 L 195 568 L 175 555 L 162 525 L 155 495 L 150 462 Z',
    labelX: 210,
    labelY: 520,
  },

  // === EUROPE === (x 405-680, y 30-260)
  // Distinctive shape with Scandinavian peninsula, Iberian bump, Italian boot area

  {
    id: 'iceland',
    // Island territory - water gap from everything. Small island northwest of Europe
    path: 'M 418 42 L 448 35 L 470 42 L 472 60 L 460 75 L 435 78 L 415 68 L 412 52 L 418 42 Z',
    labelX: 442,
    labelY: 58,
  },
  {
    id: 'scandinavia',
    // Top right of Europe, elongated north-south. Shares south with northern_europe, west gap to GB, east with ukraine
    path: 'M 500 32 L 540 28 L 575 35 L 585 60 L 578 90 L 562 115 L 540 130 L 510 132 L 495 115 L 490 85 L 492 55 L 500 32 Z',
    labelX: 535,
    labelY: 78,
  },
  {
    id: 'great_britain',
    // Island territory west of northern europe. Small island shape
    path: 'M 428 92 L 455 85 L 468 95 L 470 120 L 462 145 L 445 155 L 428 150 L 418 132 L 418 108 L 428 92 Z',
    labelX: 444,
    labelY: 120,
  },
  {
    id: 'northern_europe',
    // Center of Europe. Shares north with scandinavia, east with ukraine, south with western/southern europe
    path: 'M 485 120 L 512 134 L 542 132 L 570 120 L 582 140 L 575 165 L 555 180 L 520 185 L 490 178 L 478 158 L 478 138 L 485 120 Z',
    labelX: 525,
    labelY: 155,
  },
  {
    id: 'western_europe',
    // Southwest Europe (Iberia, France). Shares north-east with northern_europe, east with southern_europe
    path: 'M 410 155 L 440 148 L 478 140 L 480 160 L 492 180 L 488 210 L 475 238 L 450 252 L 425 248 L 410 228 L 408 195 L 410 155 Z',
    labelX: 448,
    labelY: 200,
  },
  {
    id: 'southern_europe',
    // South-center of Europe (Italy, Balkans). Shares west with western_europe, north with northern_europe, east with ukraine
    path: 'M 494 182 L 522 187 L 557 182 L 577 167 L 592 180 L 588 210 L 575 238 L 555 255 L 525 260 L 500 252 L 477 240 L 490 212 L 494 182 Z',
    labelX: 535,
    labelY: 220,
  },
  {
    id: 'ukraine',
    // Large territory, east of Europe. Shares west with scandinavia/northern_europe/southern_europe, east transitions to Asia
    path: 'M 577 37 L 620 32 L 665 38 L 678 60 L 675 95 L 668 130 L 655 165 L 635 190 L 610 200 L 590 212 L 594 182 L 584 165 L 578 140 L 584 118 L 577 90 L 580 60 L 577 37 Z',
    labelX: 628,
    labelY: 115,
  },

  // === AFRICA === (x 395-670, y 255-555)
  // Classic shield/heart shape - wide middle, narrows at south

  {
    id: 'north_africa',
    // Top-left of Africa, very wide. Shares east with egypt, south with congo/east_africa
    path: 'M 400 258 L 455 255 L 510 258 L 545 265 L 548 290 L 540 320 L 520 345 L 490 355 L 450 358 L 415 345 L 398 315 L 395 285 L 400 258 Z',
    labelX: 470,
    labelY: 305,
  },
  {
    id: 'egypt',
    // Top-right of Africa. Shares west with north_africa, south with east_africa
    path: 'M 547 265 L 585 258 L 625 260 L 650 275 L 648 305 L 635 330 L 610 345 L 580 348 L 555 340 L 542 322 L 550 290 L 547 265 Z',
    labelX: 595,
    labelY: 300,
  },
  {
    id: 'east_africa',
    // East side of Africa. Shares north with egypt, west with north_africa/congo, south with south_africa
    path: 'M 557 342 L 582 350 L 612 347 L 637 332 L 660 345 L 662 380 L 652 415 L 630 440 L 600 448 L 570 440 L 548 418 L 540 390 L 540 360 L 557 342 Z',
    labelX: 600,
    labelY: 390,
  },
  {
    id: 'congo',
    // Center-west of Africa. Shares north with north_africa, east with east_africa, south with south_africa
    path: 'M 417 347 L 452 360 L 492 357 L 522 347 L 542 362 L 542 392 L 550 420 L 535 445 L 510 455 L 478 452 L 450 440 L 428 415 L 415 385 L 410 360 L 417 347 Z',
    labelX: 478,
    labelY: 400,
  },
  {
    id: 'south_africa',
    // Bottom of Africa, narrows to point. Shares north with congo and east_africa
    path: 'M 452 442 L 480 454 L 512 457 L 537 447 L 572 442 L 602 450 L 615 475 L 610 510 L 590 535 L 560 550 L 530 555 L 498 548 L 470 530 L 450 505 L 440 475 L 445 455 L 452 442 Z',
    labelX: 530,
    labelY: 500,
  },
  {
    id: 'madagascar',
    // Island off east coast of Africa - water gap from east_africa/south_africa
    path: 'M 638 455 L 658 450 L 668 462 L 670 490 L 662 515 L 648 525 L 635 518 L 630 495 L 632 470 L 638 455 Z',
    labelX: 650,
    labelY: 488,
  },

  // === ASIA === (x 570-1065, y 10-345)
  // Massive continent spanning right side of map

  {
    id: 'ural',
    // West edge of Asia, connects to Europe (ukraine). Shares east with siberia/afghanistan, south with afghanistan
    path: 'M 680 38 L 720 30 L 755 35 L 760 60 L 758 90 L 750 120 L 735 145 L 710 150 L 685 142 L 672 118 L 670 85 L 675 55 L 680 38 Z',
    labelX: 715,
    labelY: 88,
  },
  {
    id: 'siberia',
    // Large territory across top of Asia. Shares west with ural, east with yakutsk, south with irkutsk/mongolia
    path: 'M 757 35 L 800 22 L 845 18 L 870 28 L 868 60 L 860 95 L 845 125 L 820 140 L 790 142 L 762 130 L 752 105 L 755 70 L 757 35 Z',
    labelX: 810,
    labelY: 78,
  },
  {
    id: 'yakutsk',
    // Upper right of Asia. Shares west with siberia, east with kamchatka, south with irkutsk
    path: 'M 872 28 L 915 18 L 955 22 L 965 45 L 958 75 L 942 100 L 918 112 L 890 110 L 868 95 L 862 65 L 868 42 L 872 28 Z',
    labelX: 915,
    labelY: 65,
  },
  {
    id: 'kamchatka',
    // Far east peninsula of Asia. Shares west with yakutsk, south with mongolia (indirect)
    path: 'M 957 22 L 1000 12 L 1040 15 L 1060 30 L 1062 60 L 1055 90 L 1038 112 L 1012 120 L 985 115 L 968 98 L 960 72 L 957 22 Z',
    labelX: 1012,
    labelY: 65,
  },
  {
    id: 'irkutsk',
    // Below siberia/yakutsk. Shares north with siberia/yakutsk, east with mongolia, south with mongolia/china
    path: 'M 822 142 L 860 128 L 892 112 L 920 114 L 944 102 L 952 128 L 945 158 L 925 175 L 895 178 L 862 175 L 838 165 L 822 142 Z',
    labelX: 888,
    labelY: 148,
  },
  {
    id: 'mongolia',
    // Below irkutsk, east of china. Shares with irkutsk, china, and kamchatka indirectly
    path: 'M 864 177 L 897 180 L 927 177 L 955 165 L 978 155 L 995 170 L 992 200 L 978 228 L 950 245 L 915 250 L 882 245 L 860 225 L 855 200 L 864 177 Z',
    labelX: 925,
    labelY: 212,
  },
  {
    id: 'japan',
    // Island territory east of Asia - water gap from kamchatka/mongolia/china
    path: 'M 1025 132 L 1048 125 L 1060 138 L 1062 165 L 1058 195 L 1048 215 L 1032 220 L 1020 208 L 1015 180 L 1018 155 L 1025 132 Z',
    labelX: 1040,
    labelY: 175,
  },
  {
    id: 'afghanistan',
    // Central-west Asia. Shares north with ural, east with china, south with india/middle_east
    path: 'M 687 144 L 712 152 L 737 147 L 755 125 L 768 140 L 778 165 L 775 195 L 762 222 L 735 235 L 708 230 L 685 215 L 675 190 L 672 165 L 680 148 L 687 144 Z',
    labelX: 722,
    labelY: 188,
  },
  {
    id: 'china',
    // Very large territory. Shares west with afghanistan, north with siberia/irkutsk/mongolia, east with mongolia, south with india/siam
    path: 'M 780 165 L 810 148 L 840 142 L 866 148 L 862 177 L 857 202 L 862 228 L 885 247 L 920 252 L 952 247 L 960 270 L 945 300 L 915 318 L 878 325 L 840 315 L 810 295 L 790 268 L 778 240 L 770 210 L 775 185 L 780 165 Z',
    labelX: 865,
    labelY: 240,
  },
  {
    id: 'india',
    // South-central Asia, peninsula shape. Shares north with afghanistan/china, east with siam/china
    path: 'M 710 232 L 737 237 L 764 224 L 780 242 L 792 270 L 812 297 L 808 325 L 790 345 L 762 340 L 735 330 L 715 310 L 705 285 L 700 260 L 710 232 Z',
    labelX: 752,
    labelY: 290,
  },
  {
    id: 'siam',
    // Southeast Asia. Shares north with china, west with india (indirect)
    path: 'M 842 317 L 880 327 L 917 320 L 947 302 L 960 315 L 955 340 L 935 345 L 910 342 L 885 345 L 862 340 L 845 335 L 842 317 Z',
    labelX: 898,
    labelY: 332,
  },
  {
    id: 'middle_east',
    // Southwest Asia, bridge to Africa/Europe. Shares north with ukraine(indirect)/afghanistan, east with india/afghanistan
    path: 'M 592 214 L 625 202 L 658 195 L 677 192 L 677 218 L 687 240 L 712 234 L 708 260 L 698 285 L 675 300 L 645 305 L 618 298 L 598 278 L 585 252 L 585 230 L 592 214 Z',
    labelX: 642,
    labelY: 255,
  },

  // === AUSTRALIA === (x 870-1075, y 355-565)
  // Indonesia and New Guinea as islands, western/eastern australia as mainland

  {
    id: 'indonesia',
    // Island chain - water gap from siam and australia. Elongated archipelago shape
    path: 'M 878 370 L 915 362 L 948 368 L 962 385 L 958 405 L 940 418 L 912 420 L 888 412 L 875 395 L 878 370 Z',
    labelX: 918,
    labelY: 392,
  },
  {
    id: 'new_guinea',
    // Island east of Indonesia - water gap from indonesia and eastern australia
    path: 'M 985 362 L 1025 355 L 1058 362 L 1070 380 L 1065 400 L 1048 415 L 1020 418 L 995 410 L 982 395 L 985 362 Z',
    labelX: 1025,
    labelY: 388,
  },
  {
    id: 'western_australia',
    // West half of Australian mainland. Shares east edge with eastern_australia
    path: 'M 885 440 L 920 435 L 955 438 L 972 448 L 975 480 L 975 515 L 968 542 L 948 558 L 920 562 L 895 550 L 880 525 L 875 495 L 878 462 L 885 440 Z',
    labelX: 925,
    labelY: 500,
  },
  {
    id: 'eastern_australia',
    // East half of Australian mainland. Shares west edge with western_australia
    path: 'M 977 448 L 1010 440 L 1042 445 L 1062 462 L 1072 490 L 1072 525 L 1060 550 L 1038 562 L 1010 565 L 985 555 L 977 530 L 977 498 L 977 465 L 977 448 Z',
    labelX: 1025,
    labelY: 505,
  },
];

// Connection lines between territories for visual context (dashed lines across water)
export const connectionLines: { from: string; to: string; path: string }[] = [
  // Alaska - Kamchatka (wraps around map edges)
  { from: 'alaska', to: 'kamchatka', path: 'M 27 55 Q 15 35, 5 25 M 1095 25 Q 1075 15, 1060 20' },
  // Greenland - Iceland
  { from: 'greenland', to: 'iceland', path: 'M 300 40 Q 290 42, 472 52' },
  // Greenland - NW Territory
  { from: 'greenland', to: 'northwest_territory', path: 'M 295 70 Q 285 75, 275 55' },
  // Greenland - Quebec
  { from: 'greenland', to: 'quebec', path: 'M 320 90 Q 330 95, 340 100' },
  // Iceland - Scandinavia
  { from: 'iceland', to: 'scandinavia', path: 'M 470 50 Q 480 42, 500 36' },
  // Iceland - Great Britain
  { from: 'iceland', to: 'great_britain', path: 'M 435 78 Q 432 85, 432 92' },
  // Great Britain - Northern Europe
  { from: 'great_britain', to: 'northern_europe', path: 'M 468 110 Q 476 112, 486 118' },
  // Great Britain - Western Europe
  { from: 'great_britain', to: 'western_europe', path: 'M 445 155 Q 442 158, 440 155' },
  // Brazil - North Africa
  { from: 'brazil', to: 'north_africa', path: 'M 320 400 Q 355 350, 400 290' },
  // Southern Europe - Egypt
  { from: 'southern_europe', to: 'egypt', path: 'M 560 255 Q 575 258, 585 260' },
  // Southern Europe - North Africa
  { from: 'southern_europe', to: 'north_africa', path: 'M 530 260 Q 520 262, 510 260' },
  // Western Europe - North Africa
  { from: 'western_europe', to: 'north_africa', path: 'M 430 252 Q 425 258, 420 260' },
  // East Africa - Madagascar
  { from: 'east_africa', to: 'madagascar', path: 'M 630 445 Q 635 450, 640 455' },
  // East Africa - Middle East
  { from: 'east_africa', to: 'middle_east', path: 'M 660 348 Q 665 330, 660 310' },
  // Egypt - Middle East
  { from: 'egypt', to: 'middle_east', path: 'M 650 280 Q 660 282, 670 290' },
  // Siam - Indonesia
  { from: 'siam', to: 'indonesia', path: 'M 900 345 Q 905 355, 910 365' },
  // Indonesia - New Guinea
  { from: 'indonesia', to: 'new_guinea', path: 'M 960 390 Q 970 385, 982 380' },
  // Indonesia - Western Australia
  { from: 'indonesia', to: 'western_australia', path: 'M 918 420 Q 920 428, 920 438' },
  // New Guinea - Eastern Australia
  { from: 'new_guinea', to: 'eastern_australia', path: 'M 1035 418 Q 1035 428, 1035 442' },
  // Central America - Venezuela
  { from: 'central_america', to: 'venezuela', path: 'M 170 328 Q 165 330, 155 332' },
  // Alaska - Northwest Territory
  { from: 'alaska', to: 'northwest_territory', path: 'M 135 70 L 137 70' },
];

// Continent colors for the background regions
export const continentColors: Record<string, string> = {
  north_america: 'rgba(255, 200, 50, 0.08)',
  south_america: 'rgba(233, 69, 96, 0.08)',
  europe: 'rgba(100, 149, 237, 0.08)',
  africa: 'rgba(255, 165, 0, 0.08)',
  asia: 'rgba(50, 205, 50, 0.08)',
  australia: 'rgba(148, 103, 189, 0.08)',
};
