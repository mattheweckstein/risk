export interface TerritoryMapData {
  id: string;
  path: string;
  labelX: number;
  labelY: number;
}

export const territoryPaths: TerritoryMapData[] = [
  // === NORTH AMERICA === (x 25-370, y 20-330)

  {
    id: 'alaska',
    // Alaska: jutting northwest with peninsula feel, curved coastlines
    // Shared border with northwest_territory along east edge: L 130 48 L 128 80 L 120 105
    path: 'M 38 50 C 30 38, 55 22, 80 28 C 95 30, 110 35, 130 48 L 128 80 L 120 105 C 105 115, 75 120, 55 118 C 35 112, 28 90, 30 72 C 32 60, 36 54, 38 50 Z',
    labelX: 75,
    labelY: 75,
  },
  {
    id: 'northwest_territory',
    // Wide territory across top. Shared west with alaska: M 130 48 ... L 120 105 (reversed)
    // Shared south-west with alberta: L 128 112 L 130 145
    // Shared south with ontario: L 195 140 L 240 120 L 268 95
    path: 'M 130 48 C 155 35, 190 28, 230 26 C 260 25, 280 30, 290 40 L 268 95 L 240 120 L 195 140 L 130 145 L 128 112 L 120 105 L 128 80 L 130 48 Z',
    labelX: 200,
    labelY: 80,
  },
  {
    id: 'greenland',
    // Large island northeast, curved coastlines, no shared land borders
    path: 'M 310 22 C 330 16, 360 14, 385 20 C 400 25, 408 38, 405 55 C 402 72, 392 88, 375 98 C 358 105, 335 102, 318 92 C 302 82, 295 65, 298 48 C 300 35, 305 26, 310 22 Z',
    labelX: 350,
    labelY: 58,
  },
  {
    id: 'alberta',
    // Below NW territory on west side
    // Shared north with NW territory: M 130 145 L 195 140 (reversed from NW's south)
    // Shared east with ontario: L 195 140 L 198 175 L 195 210
    // Shared south with western_us: L 145 218 L 100 210
    path: 'M 130 145 L 195 140 L 198 175 L 195 210 L 145 218 L 100 210 C 92 195, 90 170, 95 155 C 100 148, 115 145, 130 145 Z',
    labelX: 148,
    labelY: 178,
  },
  {
    id: 'ontario',
    // Center of NA
    // Shared west with alberta: M 195 140 ... L 195 210 (reversed)
    // Shared north with NW territory: M 240 120 L 268 95 (reversed from NW's south)
    // Shared east with quebec: L 290 108 L 302 140 L 295 185
    // Shared south with eastern_us: L 250 210 L 195 210
    path: 'M 195 140 L 240 120 L 268 95 L 290 108 L 302 140 L 295 185 L 250 210 L 195 210 L 198 175 L 195 140 Z',
    labelX: 245,
    labelY: 155,
  },
  {
    id: 'quebec',
    // East side below greenland
    // Shared west with ontario: M 290 108 L 302 140 L 295 185 (reversed)
    // Shared south with eastern_us: L 295 185 L 278 200
    path: 'M 290 108 C 305 90, 320 82, 340 88 C 355 94, 362 108, 358 128 C 354 148, 340 168, 320 180 L 295 185 L 302 140 L 290 108 Z',
    labelX: 325,
    labelY: 135,
  },
  {
    id: 'western_us',
    // Below alberta, west of eastern_us
    // Shared north with alberta: M 100 210 L 145 218 L 195 210 (reversed from alberta's south)
    // Shared east with eastern_us: L 195 210 L 200 245 L 185 278
    // Shared south with central_america: L 140 290 L 105 278
    path: 'M 100 210 L 145 218 L 195 210 L 200 245 L 185 278 L 140 290 L 105 278 C 80 265, 68 248, 65 232 C 64 220, 75 212, 100 210 Z',
    labelX: 132,
    labelY: 248,
  },
  {
    id: 'eastern_us',
    // East of western_us, below ontario/quebec
    // Shared west with western_us: M 195 210 L 200 245 L 185 278 (reversed)
    // Shared north with ontario: L 250 210 L 295 185 (reversed from ontario's south)
    // Shared north-east with quebec: L 295 185 L 320 180 (reversed from quebec's south)
    path: 'M 195 210 L 250 210 L 295 185 L 320 180 C 330 195, 328 215, 320 235 C 310 258, 290 275, 265 285 C 240 292, 210 290, 185 278 L 200 245 L 195 210 Z',
    labelX: 258,
    labelY: 240,
  },
  {
    id: 'central_america',
    // Narrow bottom, bridge to SA
    // Shared north with western_us: M 105 278 L 140 290 (reversed)
    // Shared north with eastern_us: L 185 278 (reversed from eastern_us)
    path: 'M 105 278 L 140 290 L 185 278 C 195 288, 200 302, 195 315 C 188 330, 172 340, 158 342 C 142 343, 128 338, 118 328 C 108 318, 100 305, 98 292 C 99 285, 102 280, 105 278 Z',
    labelX: 150,
    labelY: 312,
  },

  // === SOUTH AMERICA === (x 130-330, y 330-570)

  {
    id: 'venezuela',
    // Top of SA, wide across north
    // Shared south with peru: L 160 395 L 155 380
    // Shared south-east with brazil: L 210 405 L 260 395
    path: 'M 155 345 C 175 338, 210 335, 250 340 C 275 344, 295 352, 300 365 C 300 375, 290 385, 260 395 L 210 405 L 160 395 L 155 380 C 140 370, 138 358, 145 350 C 148 346, 152 345, 155 345 Z',
    labelX: 220,
    labelY: 368,
  },
  {
    id: 'peru',
    // West side of SA
    // Shared north with venezuela: M 155 380 L 160 395 (reversed)
    // Shared east with brazil: L 210 405 L 215 440 L 205 470
    // Shared south with argentina: L 175 480 L 150 470
    path: 'M 155 380 L 160 395 L 210 405 L 215 440 L 205 470 L 175 480 L 150 470 C 135 455, 128 435, 130 415 C 132 400, 140 390, 155 380 Z',
    labelX: 172,
    labelY: 432,
  },
  {
    id: 'brazil',
    // Largest in SA, bulges east
    // Shared west with venezuela: M 260 395 L 210 405 (reversed)
    // Shared west with peru: M 210 405 L 215 440 L 205 470 (reversed)
    // Shared south with argentina: L 250 492 L 205 470
    path: 'M 260 395 C 280 388, 300 385, 318 392 C 332 400, 335 418, 330 440 C 325 460, 310 478, 290 488 C 270 495, 250 495, 250 492 L 205 470 L 215 440 L 210 405 L 260 395 Z',
    labelX: 275,
    labelY: 440,
  },
  {
    id: 'argentina',
    // Bottom of SA, narrows south
    // Shared north with peru: M 150 470 L 175 480 (reversed)
    // Shared north with brazil: M 205 470 L 250 492 (reversed)
    path: 'M 150 470 L 175 480 L 205 470 L 250 492 C 258 510, 252 530, 240 548 C 228 562, 212 572, 198 575 C 182 575, 168 565, 160 548 C 152 530, 148 510, 145 492 C 144 480, 146 474, 150 470 Z',
    labelX: 200,
    labelY: 525,
  },

  // === EUROPE === (x 400-680, y 30-260)

  {
    id: 'iceland',
    // Small island, no shared land borders
    path: 'M 420 45 C 430 38, 450 35, 465 40 C 475 44, 478 54, 474 65 C 470 74, 458 80, 445 80 C 432 79, 420 72, 416 62 C 413 54, 415 48, 420 45 Z',
    labelX: 446,
    labelY: 60,
  },
  {
    id: 'scandinavia',
    // Elongated peninsula reaching north
    // Shared south with northern_europe: L 530 145 L 505 142
    // Shared east with ukraine: L 585 52 L 588 90 L 580 128
    path: 'M 510 32 C 525 26, 550 25, 570 30 L 585 52 L 588 90 L 580 128 L 530 145 L 505 142 C 495 130, 492 110, 494 90 C 496 70, 500 50, 510 32 Z',
    labelX: 538,
    labelY: 85,
  },
  {
    id: 'great_britain',
    // Island, no shared land borders
    path: 'M 432 95 C 442 88, 456 87, 466 92 C 474 98, 476 110, 474 125 C 472 138, 462 148, 450 152 C 438 155, 426 148, 422 136 C 418 122, 420 108, 425 100 C 428 96, 430 95, 432 95 Z',
    labelX: 448,
    labelY: 122,
  },
  {
    id: 'northern_europe',
    // Center of Europe
    // Shared north with scandinavia: M 505 142 L 530 145 (reversed)
    // Shared east with ukraine: L 580 128 L 585 162 (uses shared scandinavia-ukraine border then extends)
    // Shared south with southern_europe: L 560 188 L 525 195 L 498 190
    // Shared south-west with western_europe: L 485 170 L 480 150
    path: 'M 505 142 L 530 145 L 580 128 L 585 162 L 560 188 L 525 195 L 498 190 L 485 170 L 480 150 L 505 142 Z',
    labelX: 530,
    labelY: 165,
  },
  {
    id: 'western_europe',
    // Southwest Europe (Iberia, France)
    // Shared east with northern_europe: M 480 150 L 485 170 (reversed)
    // Shared east with southern_europe: L 498 190 L 495 220 L 482 248
    path: 'M 480 150 L 485 170 L 498 190 L 495 220 L 482 248 C 465 258, 440 260, 425 252 C 412 242, 408 225, 408 205 C 408 185, 412 168, 420 158 C 435 148, 458 148, 480 150 Z',
    labelX: 448,
    labelY: 208,
  },
  {
    id: 'southern_europe',
    // South-center (Italy, Balkans)
    // Shared north with northern_europe: M 498 190 L 525 195 L 560 188 (reversed)
    // Shared east with ukraine: L 585 162 L 600 195 L 598 225
    // Shared west with western_europe: M 498 190 L 495 220 L 482 248 (reversed)
    path: 'M 498 190 L 525 195 L 560 188 L 585 162 L 600 195 L 598 225 C 592 242, 578 255, 560 262 C 540 268, 518 265, 500 258 C 488 252, 483 248, 482 248 L 495 220 L 498 190 Z',
    labelX: 540,
    labelY: 228,
  },
  {
    id: 'ukraine',
    // Large territory, east of Europe bridging to Asia
    // Shared west with scandinavia: M 585 52 L 588 90 L 580 128 (reversed)
    // Shared south-west with northern_europe: M 580 128 L 585 162 (reversed)
    // Shared south with southern_europe: M 585 162 L 600 195 L 598 225 (reversed)
    path: 'M 585 52 C 610 38, 645 35, 672 42 C 690 48, 695 62, 692 82 C 688 105, 678 135, 665 165 C 655 188, 640 205, 620 218 L 598 225 L 600 195 L 585 162 L 580 128 L 588 90 L 585 52 Z',
    labelX: 635,
    labelY: 125,
  },

  // === AFRICA === (x 390-670, y 255-560)

  {
    id: 'north_africa',
    // Top-left of Africa, wide
    // Shared east with egypt: L 555 272 L 555 310 L 548 340
    // Shared south-east with east_africa: L 548 340 L 535 365
    // Shared south with congo: L 500 372 L 450 370 L 415 358
    path: 'M 402 262 C 430 256, 480 255, 520 260 L 555 272 L 555 310 L 548 340 L 535 365 L 500 372 L 450 370 L 415 358 C 400 345, 392 325, 392 305 C 392 285, 395 270, 402 262 Z',
    labelX: 472,
    labelY: 312,
  },
  {
    id: 'egypt',
    // Top-right of Africa
    // Shared west with north_africa: M 555 272 L 555 310 L 548 340 (reversed)
    // Shared south with east_africa: L 548 340 L 565 348 L 600 348
    path: 'M 555 272 C 575 262, 605 258, 630 265 C 645 270, 655 282, 655 298 C 654 318, 640 338, 620 348 L 600 348 L 565 348 L 548 340 L 555 310 L 555 272 Z',
    labelX: 598,
    labelY: 305,
  },
  {
    id: 'east_africa',
    // East side
    // Shared north with egypt: M 565 348 L 600 348 L 620 348 (reversed from egypt's south)
    // Shared north-west with north_africa: M 535 365 L 548 340 (reversed)
    // Shared west with congo: L 535 365 L 540 400 L 548 430
    // Shared south with south_africa: L 602 455 L 580 452
    path: 'M 548 340 L 565 348 L 600 348 L 620 348 C 645 355, 660 375, 665 400 C 668 425, 660 448, 645 462 C 630 472, 615 465, 602 455 L 580 452 L 548 430 L 540 400 L 535 365 L 548 340 Z',
    labelX: 600,
    labelY: 400,
  },
  {
    id: 'congo',
    // Center-west
    // Shared north with north_africa: M 415 358 L 450 370 L 500 372 L 535 365 (reversed)
    // Shared east with east_africa: M 535 365 L 540 400 L 548 430 (reversed)
    // Shared south with south_africa: L 548 430 L 518 462 L 470 458
    path: 'M 415 358 L 450 370 L 500 372 L 535 365 L 540 400 L 548 430 L 518 462 L 470 458 C 445 452, 425 438, 415 420 C 405 402, 402 382, 408 365 C 410 360, 413 358, 415 358 Z',
    labelX: 478,
    labelY: 410,
  },
  {
    id: 'south_africa',
    // Bottom, narrows to point
    // Shared north with congo: M 470 458 L 518 462 (reversed)
    // Shared north-east with east_africa: M 548 430 L 580 452 L 602 455 (reversed)
    path: 'M 470 458 L 518 462 L 548 430 L 580 452 L 602 455 C 618 472, 622 498, 615 520 C 608 542, 590 558, 568 565 C 545 570, 520 568, 500 558 C 478 548, 462 530, 455 508 C 450 490, 452 472, 460 462 C 464 458, 468 458, 470 458 Z',
    labelX: 535,
    labelY: 510,
  },
  {
    id: 'madagascar',
    // Island off southeast coast, no shared land borders
    path: 'M 645 462 C 655 455, 668 455, 675 465 C 680 475, 680 492, 675 510 C 670 525, 660 535, 648 532 C 638 528, 632 515, 632 498 C 632 480, 636 468, 645 462 Z',
    labelX: 655,
    labelY: 498,
  },

  // === ASIA === (x 570-1065, y 10-350)

  {
    id: 'ural',
    // West edge of Asia
    // Shared west (connects to ukraine via adjacency, not shared border - separated by map gap)
    // Shared east with siberia: L 755 42 L 758 82 L 752 120
    // Shared south with afghanistan: L 720 155 L 690 148
    // Shared south-east with china: L 752 120 L 755 145 (through siberia connection)
    path: 'M 695 38 C 710 30, 735 28, 755 42 L 758 82 L 752 120 L 720 155 L 690 148 C 678 138, 675 118, 678 98 C 680 75, 685 52, 695 38 Z',
    labelX: 718,
    labelY: 92,
  },
  {
    id: 'siberia',
    // Large across top of Asia
    // Shared west with ural: M 755 42 L 758 82 L 752 120 (reversed)
    // Shared east with yakutsk: L 870 32 L 868 72 L 862 108
    // Shared south with irkutsk: L 840 138 L 808 148
    // Shared south-west with china: (through irkutsk)
    path: 'M 755 42 C 780 28, 820 20, 850 22 L 870 32 L 868 72 L 862 108 L 840 138 L 808 148 C 785 148, 765 140, 752 120 L 758 82 L 755 42 Z',
    labelX: 810,
    labelY: 82,
  },
  {
    id: 'yakutsk',
    // Upper right
    // Shared west with siberia: M 870 32 L 868 72 L 862 108 (reversed)
    // Shared east with kamchatka: L 960 28 L 962 68 L 955 105
    // Shared south with irkutsk: L 862 108 L 895 118
    path: 'M 870 32 C 895 22, 930 18, 960 28 L 962 68 L 955 105 L 895 118 L 862 108 L 868 72 L 870 32 Z',
    labelX: 915,
    labelY: 68,
  },
  {
    id: 'kamchatka',
    // Far east peninsula
    // Shared west with yakutsk: M 960 28 L 962 68 L 955 105 (reversed)
    // Shared south-west with mongolia: L 955 105 L 968 120
    path: 'M 960 28 C 985 18, 1020 12, 1048 18 C 1065 24, 1070 42, 1068 62 C 1065 85, 1055 105, 1038 118 C 1020 128, 998 130, 978 125 L 968 120 L 955 105 L 962 68 L 960 28 Z',
    labelX: 1018,
    labelY: 68,
  },
  {
    id: 'irkutsk',
    // Below siberia/yakutsk
    // Shared north with siberia: M 808 148 L 840 138 (reversed)
    // Shared north-east with yakutsk: M 862 108 L 895 118 (reversed)
    // Shared south with mongolia: L 948 155 L 940 180 L 910 195
    // Shared south-west with china: L 870 192 L 838 180
    path: 'M 808 148 L 840 138 L 862 108 L 895 118 L 948 155 L 940 180 L 910 195 L 870 192 L 838 180 L 808 148 Z',
    labelX: 880,
    labelY: 162,
  },
  {
    id: 'mongolia',
    // Below irkutsk
    // Shared north with irkutsk: M 910 195 L 940 180 L 948 155 (reversed)
    // Shared north-east with kamchatka: M 968 120 (reversed from kamchatka's south)
    // Shared south with china: L 925 262 L 880 258 L 858 238
    path: 'M 910 195 L 940 180 L 948 155 L 968 120 L 978 125 C 995 140, 1002 162, 1000 185 C 998 210, 985 235, 965 252 C 948 265, 930 268, 925 262 L 880 258 L 858 238 C 858 218, 868 205, 882 198 L 910 195 Z',
    labelX: 930,
    labelY: 218,
  },
  {
    id: 'japan',
    // Island territory, no shared land borders
    path: 'M 1030 138 C 1042 132, 1055 132, 1062 142 C 1068 152, 1068 170, 1065 190 C 1062 208, 1052 222, 1040 225 C 1028 227, 1020 218, 1018 205 C 1015 188, 1018 165, 1022 150 C 1025 142, 1028 138, 1030 138 Z',
    labelX: 1042,
    labelY: 180,
  },
  {
    id: 'afghanistan',
    // Central-west Asia
    // Shared north with ural: M 690 148 L 720 155 (reversed)
    // Shared east with china: L 778 175 L 775 210 L 762 240
    // Shared south with india: L 730 248 L 705 238
    // Shared south-west with middle_east: L 682 218 L 680 185
    path: 'M 690 148 L 720 155 L 752 120 L 808 148 L 838 180 L 778 175 L 775 210 L 762 240 L 730 248 L 705 238 L 682 218 L 680 185 L 690 148 Z',
    labelX: 730,
    labelY: 195,
  },
  {
    id: 'china',
    // Very large territory
    // Shared west with afghanistan: M 778 175 L 775 210 L 762 240 (reversed)
    // Shared north-west with siberia/irkutsk border: M 838 180 (shared)
    // Shared north with irkutsk: M 838 180 L 870 192 (reversed)
    // Shared north-east with mongolia: M 858 238 L 880 258 L 925 262 (reversed)
    // Shared south with siam: L 912 330 L 862 338
    // Shared south-west with india: L 810 310 L 792 278
    path: 'M 778 175 L 838 180 L 870 192 L 858 238 L 880 258 L 925 262 C 948 268, 960 280, 962 298 C 962 315, 950 328, 932 335 L 912 330 L 862 338 L 810 310 L 792 278 L 762 240 L 775 210 L 778 175 Z',
    labelX: 868,
    labelY: 268,
  },
  {
    id: 'india',
    // Peninsula shape
    // Shared north with afghanistan: M 705 238 L 730 248 (reversed)
    // Shared north-east with china: M 762 240 L 792 278 L 810 310 (reversed)
    // Shared east with siam (indirect, via china connection)
    // Shared west with middle_east (indirect)
    path: 'M 705 238 L 730 248 L 762 240 L 792 278 L 810 310 C 812 328, 805 345, 790 355 C 775 362, 755 358, 740 348 C 722 335, 710 315, 705 295 C 700 275, 698 258, 702 245 L 705 238 Z',
    labelX: 752,
    labelY: 305,
  },
  {
    id: 'siam',
    // Southeast Asia peninsula
    // Shared north with china: M 862 338 L 912 330 (reversed)
    path: 'M 862 338 L 912 330 L 932 335 C 950 340, 958 348, 955 358 C 950 368, 935 372, 918 370 C 900 368, 882 362, 868 355 C 856 348, 855 342, 862 338 Z',
    labelX: 905,
    labelY: 352,
  },
  {
    id: 'middle_east',
    // Southwest Asia, bridge area
    // Shared north with ukraine (adjacency, no shared pixel border)
    // Shared north-east with afghanistan: M 680 185 L 682 218 (reversed)
    // Shared east with india: M 705 238 (reversed from india, indirect)
    path: 'M 620 218 C 635 208, 658 198, 680 185 L 682 218 L 705 238 L 702 245 C 700 262, 695 280, 685 298 C 672 315, 652 322, 632 318 C 615 312, 600 298, 592 280 C 585 262, 588 242, 598 228 C 605 222, 612 220, 620 218 Z',
    labelX: 648,
    labelY: 268,
  },

  // === AUSTRALIA === (x 870-1075, y 355-565)

  {
    id: 'indonesia',
    // Archipelago island chain, no shared land borders
    path: 'M 882 378 C 895 370, 918 365, 940 370 C 955 374, 965 385, 962 398 C 958 410, 942 418, 922 420 C 902 420, 885 414, 878 402 C 874 392, 876 384, 882 378 Z',
    labelX: 920,
    labelY: 395,
  },
  {
    id: 'new_guinea',
    // Island east of Indonesia, no shared land borders
    path: 'M 990 368 C 1005 360, 1028 358, 1048 365 C 1062 370, 1070 382, 1068 398 C 1065 412, 1052 422, 1035 425 C 1018 426, 1000 420, 990 408 C 982 398, 982 382, 990 368 Z',
    labelX: 1028,
    labelY: 395,
  },
  {
    id: 'western_australia',
    // West half of mainland
    // Shared east with eastern_australia: L 975 455 L 975 505 L 975 548
    path: 'M 888 445 C 910 436, 940 434, 960 438 L 975 455 L 975 505 L 975 548 C 965 558, 945 565, 925 568 C 905 568, 888 560, 878 545 C 868 528, 865 505, 868 482 C 872 462, 878 450, 888 445 Z',
    labelX: 925,
    labelY: 502,
  },
  {
    id: 'eastern_australia',
    // East half of mainland
    // Shared west with western_australia: M 975 455 L 975 505 L 975 548 (reversed)
    path: 'M 975 455 L 960 438 C 980 434, 1005 436, 1025 445 C 1045 455, 1058 472, 1065 495 C 1070 518, 1068 542, 1058 558 C 1045 572, 1025 575, 1005 572 C 988 568, 978 558, 975 548 L 975 505 L 975 455 Z',
    labelX: 1022,
    labelY: 510,
  },
];

// Connection lines between territories for visual context (dashed lines across water)
export const connectionLines: { from: string; to: string; path: string }[] = [
  // Alaska - Kamchatka (wraps around map edges)
  { from: 'alaska', to: 'kamchatka', path: 'M 30 50 Q 15 30, 5 20 M 1090 20 Q 1075 15, 1060 25' },
  // Greenland - Iceland
  { from: 'greenland', to: 'iceland', path: 'M 310 40 Q 295 42, 475 52' },
  // Greenland - NW Territory
  { from: 'greenland', to: 'northwest_territory', path: 'M 300 65 Q 295 70, 290 60' },
  // Greenland - Quebec
  { from: 'greenland', to: 'quebec', path: 'M 330 98 Q 338 102, 345 108' },
  // Iceland - Scandinavia
  { from: 'iceland', to: 'scandinavia', path: 'M 474 50 Q 485 42, 510 36' },
  // Iceland - Great Britain
  { from: 'iceland', to: 'great_britain', path: 'M 440 78 Q 438 86, 436 92' },
  // Great Britain - Northern Europe
  { from: 'great_britain', to: 'northern_europe', path: 'M 470 115 Q 478 118, 486 125' },
  // Great Britain - Western Europe
  { from: 'great_britain', to: 'western_europe', path: 'M 450 152 Q 448 155, 445 158' },
  // Brazil - North Africa
  { from: 'brazil', to: 'north_africa', path: 'M 330 410 Q 360 370, 402 300' },
  // Southern Europe - Egypt
  { from: 'southern_europe', to: 'egypt', path: 'M 565 260 Q 578 262, 590 265' },
  // Southern Europe - North Africa
  { from: 'southern_europe', to: 'north_africa', path: 'M 535 265 Q 525 262, 515 262' },
  // Western Europe - North Africa
  { from: 'western_europe', to: 'north_africa', path: 'M 425 255 Q 420 260, 412 262' },
  // East Africa - Madagascar
  { from: 'east_africa', to: 'madagascar', path: 'M 645 460 Q 648 462, 650 465' },
  // East Africa - Middle East
  { from: 'east_africa', to: 'middle_east', path: 'M 658 358 Q 662 340, 655 320' },
  // Egypt - Middle East
  { from: 'egypt', to: 'middle_east', path: 'M 652 290 Q 660 285, 665 280' },
  // Siam - Indonesia
  { from: 'siam', to: 'indonesia', path: 'M 905 370 Q 908 375, 912 380' },
  // Indonesia - New Guinea
  { from: 'indonesia', to: 'new_guinea', path: 'M 960 400 Q 972 395, 985 390' },
  // Indonesia - Western Australia
  { from: 'indonesia', to: 'western_australia', path: 'M 920 420 Q 922 430, 922 440' },
  // New Guinea - Eastern Australia
  { from: 'new_guinea', to: 'eastern_australia', path: 'M 1038 425 Q 1040 435, 1040 448' },
  // Central America - Venezuela
  { from: 'central_america', to: 'venezuela', path: 'M 160 342 Q 158 344, 156 346' },
  // Alaska - Northwest Territory (shared land border, but add for clarity)
  { from: 'alaska', to: 'northwest_territory', path: 'M 130 48 L 130 48' },
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

// Continent territory groupings for background hulls
export const continentTerritories: Record<string, string[]> = {
  north_america: ['alaska', 'northwest_territory', 'greenland', 'alberta', 'ontario', 'quebec', 'western_us', 'eastern_us', 'central_america'],
  south_america: ['venezuela', 'peru', 'brazil', 'argentina'],
  europe: ['iceland', 'scandinavia', 'great_britain', 'northern_europe', 'western_europe', 'southern_europe', 'ukraine'],
  africa: ['north_africa', 'egypt', 'east_africa', 'congo', 'south_africa', 'madagascar'],
  asia: ['ural', 'siberia', 'yakutsk', 'kamchatka', 'irkutsk', 'mongolia', 'japan', 'afghanistan', 'china', 'india', 'siam', 'middle_east'],
  australia: ['indonesia', 'new_guinea', 'western_australia', 'eastern_australia'],
};

// Hand-drawn continent background outlines (generous padding around territories)
export const continentOutlines: Record<string, string> = {
  north_america: 'M 20 18 C 50 10, 180 8, 300 15 C 340 18, 395 20, 400 45 C 405 75, 395 95, 370 100 C 355 105, 340 115, 335 140 C 330 165, 325 185, 335 195 C 340 200, 330 220, 310 240 C 290 260, 260 290, 220 300 C 190 308, 155 345, 135 350 C 112 355, 90 345, 85 325 C 80 305, 90 290, 95 275 C 100 260, 55 255, 42 235 C 30 215, 25 180, 28 145 C 30 120, 20 85, 22 55 C 24 35, 20 25, 20 18 Z',
  south_america: 'M 125 332 C 165 325, 230 325, 295 335 C 315 340, 338 360, 342 390 C 345 420, 340 455, 325 485 C 310 515, 280 545, 255 565 C 230 580, 195 585, 170 575 C 148 565, 135 540, 130 510 C 125 480, 118 450, 120 420 C 122 390, 118 360, 125 332 Z',
  europe: 'M 405 30 C 430 22, 480 20, 520 22 C 560 24, 600 28, 680 35 C 705 50, 705 80, 698 110 C 695 140, 680 175, 660 210 C 640 235, 615 255, 580 270 C 545 280, 505 275, 470 265 C 440 260, 400 248, 395 225 C 390 200, 395 170, 400 145 C 405 120, 400 95, 398 72 C 396 55, 400 38, 405 30 Z',
  africa: 'M 385 250 C 430 242, 510 242, 570 248 C 620 254, 665 268, 678 300 C 688 330, 680 365, 675 400 C 672 430, 680 465, 685 490 C 688 510, 680 535, 660 545 C 640 555, 615 565, 580 575 C 545 582, 505 580, 475 570 C 445 560, 425 540, 415 510 C 405 480, 395 445, 390 410 C 385 375, 380 340, 378 310 C 376 280, 380 258, 385 250 Z',
  asia: 'M 668 28 C 720 18, 810 8, 910 8 C 980 8, 1050 5, 1080 22 C 1085 40, 1082 80, 1075 120 C 1068 155, 1055 200, 1020 235 C 1000 255, 975 295, 968 330 C 965 350, 960 375, 945 385 C 925 395, 870 380, 845 365 C 820 350, 795 365, 770 370 C 740 375, 698 365, 675 340 C 655 320, 635 295, 620 275 C 600 250, 582 240, 580 225 C 575 210, 595 195, 618 230 C 640 215, 660 190, 672 160 C 682 135, 688 105, 690 80 C 692 55, 680 38, 668 28 Z',
  australia: 'M 868 358 C 910 348, 975 345, 1040 350 C 1070 355, 1082 375, 1080 405 C 1078 430, 1080 458, 1078 490 C 1076 525, 1075 555, 1065 575 C 1048 590, 1010 590, 975 585 C 940 580, 905 580, 878 568 C 858 558, 850 535, 852 505 C 854 475, 855 445, 860 418 C 862 395, 860 372, 868 358 Z',
};
