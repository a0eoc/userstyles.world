export interface UserSettings {
    colorScheme: 'dark' | 'light' | 'system';
    entriesBehavior: 'hover' | 'click' | 'no-hide';
    redirect: string;
}

const DEFAULT_SETTINGS: UserSettings = {
    colorScheme: 'system',
    entriesBehavior: 'click',
    redirect: '',
};

const localStorageKey = 'user-preferences';
let settings: UserSettings = null;

export function getSettings(): UserSettings {
    if (settings) {
        return settings;
    }
    const storedSettings = localStorage.getItem(localStorageKey);
    if (!storedSettings) {
        localStorage.setItem(localStorageKey, JSON.stringify(DEFAULT_SETTINGS));
        return DEFAULT_SETTINGS;
    }
    const savedSettings = getValidatedObject(JSON.parse(storedSettings), DEFAULT_SETTINGS);

    // Data migration, just to be sure if any new setting are added.
    // We should include the default value and save it.
    settings = {...DEFAULT_SETTINGS, ...savedSettings};
    localStorage.setItem(localStorageKey, JSON.stringify(settings));

    return settings;
}

export function storeNewSettings(newSettings: Partial<UserSettings>) {
    settings = {...settings, ...newSettings};
    localStorage.setItem(localStorageKey, JSON.stringify(settings));
}

// A niece function to make sure that the return object will only return
// key's that the compare object have.
function getValidatedObject<T>(source: any, compare: T): Partial<T> {
    const result = {};
    if (null == source || 'object' !== typeof source || Array.isArray(source)) {
        return null;
    }
    Object.keys(source).forEach((key) => {
        const value = source[key];
        if (null == value || null == compare[key]) {
            return;
        }
        if ('object' === typeof value && 'object' === typeof compare[key]) {
            result[key] = getValidatedObject(value, compare[key]);
        } else if (typeof value === typeof compare[key] && compare.hasOwnProperty(key)) {
            result[key] = value;
        }
    });
    return result;
}
