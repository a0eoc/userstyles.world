import {checkRedirect} from './page/account';
import {saveRedirect} from './page/login';
import {changeEntriesBehavior} from './page/modlog';
import {initalizeOrUpdateColorScheme as initalizeOrUpdateColorScheme} from './color-scheme';
import {broadcastReady} from './third-party';
import {saveUserSettingsButton, setValues} from './user-settings';
import {doDomOperation} from './utils/dom';
import type {UserSettings} from './utils/storage';
import {getSettings} from './utils/storage';
import {initViewStyle} from './page/view-style';
import {page404} from './page/404';

// Once settings update we should reinstalize any functionallity.
// That relies on this settings.
const onSettingsUpdate = () => initalizeOrUpdateColorScheme(getSettings().colorScheme);

const whenDOMReady = () => {
    broadcastReady();
    saveUserSettingsButton(onSettingsUpdate);
    setValues(getSettings());
};

// WhenDOMReady contains code that only should be handle
// when the DOM is ready to go.
// Any other code shouldn't depend on this setup function.
doDomOperation(whenDOMReady);

// Initalize functions that requires settings and don't depend on the DOM.
// Note that we don't save getSettings() result, as this initalize is a 1 time thing
// And having it sit in the memory is kinda useless.
initalizeOrUpdateColorScheme(getSettings().colorScheme);

const styleViewRegex = /\/style\/\d+\/(?!promote)\S*/;
function pageSpecificFunctions(settings: UserSettings) {
    switch (location.pathname) {
        case '/modlog':
            changeEntriesBehavior(settings.entriesBehavior);
            break;
        case '/login':
            saveRedirect();
            break;
        case '/account':
        case '/api/oauth/style/link':
            checkRedirect(settings.redirect);
            break;
        default:
            page404();
    }
    if (location.pathname.startsWith('/style/') && styleViewRegex.test(location.pathname)) {
        initViewStyle();
    }
}

pageSpecificFunctions(getSettings());
