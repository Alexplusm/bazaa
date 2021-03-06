# TODO: setup venv
# TODO: modules/packages

from datetime import timedelta

import sys
sys.path.append('..')

from tests.ext_system.ext_system import create_ext_system
from tests.game.game import (
    create_game, update_game, prepare_game, get_games
)
from tests.screenshot.high_load import run_high_load
from tests.utils import get_timestamp
from tests.config import get_url_start
from tests.statistics.statistics import (
    get_statistics_user
)
from tests.screenshot.screenshot import (get_screenshot_results)


# prepare game

ext_system = {
    # "extSystemId": "ext-id-5",
    "description": "some description",
    "postResultsUrl": "https://abc/lol.php"
}

game = {
    "extSystemId": None,
    "name": "new game",
    "answerType": 2,
    "startDate":   str(get_timestamp(timedelta(days=2))),
    "endDate":     str(get_timestamp(timedelta(days=3))),
    "question": "Choose answer",
    "options": "yep, nope"
}


def complete_test():
    print("Config: ", get_url_start())
    ext_system_id = create_ext_system(ext_system)
    if ext_system_id is None:
        ext_system_id = ext_system["extSystemId"]

    game["extSystemId"] = ext_system_id
    game_id = create_game(game)
    if game_id is None:
        print("error while game creation")
    update_game(game_id)
    prepare_game(game_id)

    # print("game_id: ", game_id)
    # print("ext_system_id: ", ext_system_id)

    run_high_load(ext_system_id, game_id)


def complete_test_without_ext_sys(ext_system_id):
    print("Config: ", get_url_start())
    game["extSystemId"] = ext_system_id
    game_id = create_game(game)
    if game_id is None:
        print("error while game creation")

    update_game(game_id)
    prepare_game(game_id)

    run_high_load(ext_system_id, game_id)


def test_user_statistics():
    user_id = "i-user-2"
    game_id = "5c7713c7-3960-4c0d-ae7f-c27417ed234d"
    ext_system_id = "b0e4c252-9b72-4761-b574-fff694965dcf"
    params = {
        "extSystemId": ext_system_id,
        "gameIds": game_id,
        # "totalOnly": "TrUe",
        "from": str(get_timestamp(timedelta(days=-10))),
        "to": str(get_timestamp(timedelta(days=1))),
    }
    res = get_statistics_user(user_id, params)
    print("result: ", res)


def screenshot_results():
    game_id = "79771e18-e8ec-4828-8d17-3a40644dd2c4"
    screenshot_id = "750ac32c-306b-4a67-bfd3-99979daed189"
    res = get_screenshot_results(game_id, screenshot_id)
    print(res)


def main():
    # complete_test()

    # test_user_statistics()
    # screenshot_results()

    # complete_test_without_ext_sys("0af61dfe-302e-4e00-a3cf-90be4619c9a2")

    game_id = "510c1dfa-0981-40aa-88cb-e8309556d9ec"
    ext_sys_id = "clean1"
    run_high_load(ext_sys_id, game_id)

    # ext_system_id = "ext-id-3"
    # get_games(ext_system_id)


if __name__ == "__main__":
    main()
