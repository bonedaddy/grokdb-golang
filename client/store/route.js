const qs = require('qs');
const page = require('page');

const {paths} = require('store/constants');
const {generateSlug} = require('store/utils');

const transforms = {

    validDeckSlug(state, options) {

        let {slug, deck, deckID} = options;
        ({deck, deckID} = resolveDeck(state, deck, deckID));

        const slugged = generateSlug(deck.get('name'), deckID);

        options.validDeckSlug = slugged == slug;

        return options;
    },

    toDeck(state, options = {}) {

        let {deck, deckID} = options;
        ({deck, deckID} = resolveDeck(state, deck, deckID));

        const slugged = generateSlug(deck.get('name'), deckID);
        page(`/deck/${deckID}/${slugged}`);
    },

    redirectToDeck(state, options = {}) {

        let {deck, deckID} = options;
        ({deck, deckID} = resolveDeck(state, deck, deckID));

        const slugged = generateSlug(deck.get('name'), deckID);
        page.redirect(`/deck/${deckID}/${slugged}`);
    },

    toDeckEdit(state, options = {}) {

        let {deck, deckID} = options;

        ({deck, deckID} = resolveDeck(state, deck, deckID));

        const slugged = generateSlug(deck.get('name'), deckID);
        page(`/deck/${deckID}/${slugged}/edit`);
    },

    toDeckCards(state, options = {}) {

        let {deck, deckID} = options;
        const {pageNum: _pageNum = 1, order: _order, sort: _sort} = options;

        ({deck, deckID} = resolveDeck(state, deck, deckID));

        const slugged = generateSlug(deck.get('name'), deckID);

        const params = qs.stringify({page: _pageNum, order: _order, sort: _sort});

        page(`/deck/${deckID}/${slugged}/cards?${params}`);
    },

    redirectToDeckCards(state, options = {}) {

        let {deck, deckID} = options;
        const {pageNum: _pageNum = 1, order: _order, sort: _sort} = options;

        ({deck, deckID} = resolveDeck(state, deck, deckID));

        const slugged = generateSlug(deck.get('name'), deckID);

        const params = qs.stringify({page: _pageNum, order: _order, sort: _sort});

        page.redirect(`/deck/${deckID}/${slugged}/cards?${params}`);
    },

    toDeckCardsNew(state, options = {}) {

        let {deck, deckID} = options;

        ({deck, deckID} = resolveDeck(state, deck, deckID));

        const slugged = generateSlug(deck.get('name'), deckID);
        page(`/deck/${deckID}/${slugged}/cards/new`);
    },

    toCardProfile(state, options = {}) {

        let {card, cardID} = options;

        ({card, cardID} = resolveCard(state, card, cardID));

        page(`/card/${cardID}`);
    },

    toCardProfileEdit(state, options = {}) {

        let {card, cardID} = options;

        ({card, cardID} = resolveCard(state, card, cardID));

        page(`/card/${cardID}/edit`);
    },

    toDeckReview(state, options = {}) {

        let {deck, deckID} = options;

        ({deck, deckID} = resolveDeck(state, deck, deckID));

        page(`/review/deck/${deckID}`);
    },

    toDeckReviewEdit(state, options = {}) {

        let {deck, deckID} = options;

        ({deck, deckID} = resolveDeck(state, deck, deckID));

        page(`/review/deck/${deckID}/edit`);
    },

    toStash(state, options = {}) {
        page(`/stashes`);
        return options;
    },

    toStashNew(state, options = {}) {
        page(`/newstash`);
        return options;
    },

    toStashProfile(state, options = {}) {

        let {stash, stashID} = options;
        const {page: _pageNum = 1, order: _order, sort: _sort} = options;

        ({stash, stashID} = resolveStash(state, stash, stashID));

        const params = qs.stringify({page: _pageNum, order: _order, sort: _sort});

        page(`/stashes/${stashID}?${params}`);
    },

    toStashProfileEdit(state, options = {}) {

        let {stash, stashID} = options;
        const {page: _pageNum = 1, order: _order, sort: _sort} = options;

        ({stash, stashID} = resolveStash(state, stash, stashID));

        const params = qs.stringify({page: _pageNum, order: _order, sort: _sort});

        page(`/stashes/${stashID}/edit?${params}`);
    },

    toStashProfileReview(state, options = {}) {

        let {stash, stashID} = options;

        ({stash, stashID} = resolveStash(state, stash, stashID));

        page(`/stashes/${stashID}/review`);
    },

    toStashProfileReviewEdit(state, options = {}) {

        let {stash, stashID} = options;

        ({stash, stashID} = resolveStash(state, stash, stashID));

        page(`/stashes/${stashID}/review/edit`);
    }
};

module.exports = transforms;

/* helpers */

const resolveDeck = function(state, deck = void 0, deckID = void 0) {
    state.cursor(paths.deck.self).update(function(_deck) {
        deck = deck || _deck;
        return deck;
    });

    if(!deckID) {
        deckID = deck.get('id');
    }

    return {deck, deckID};
};

const resolveCard = function(state, card = void 0, cardID = void 0) {
    state.cursor(paths.card.self).update(function(_card) {
        card = card || _card;
        return card;
    });

    if(!cardID) {
        cardID = card.get('id');
    }

    return {card, cardID};
};

const resolveStash = function(state, stash = void 0, stashID = void 0) {
    state.cursor(paths.stash.self).update(function(_stash) {
        stash = stash || _stash;
        return stash;
    });

    if(!stashID) {
        stashID = stash.get('id');
    }

    return {stash, stashID};
};
