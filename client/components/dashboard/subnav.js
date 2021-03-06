const React = require('react');
const orwell = require('orwell');
const classNames = require('classnames');

const {flow} = require('store/utils');
const {dashboard, paths} = require('store/constants');
const {toDeckCards, toDeck, toDeckReview, toStash} = require('store/route');
const {applyDeckCardsPageArgs} = require('store/cards');

const toDeckCardsList = flow(
    applyDeckCardsPageArgs,
    toDeckCards
);

const SubNav = React.createClass({

    propTypes: {
        store: React.PropTypes.object.isRequired,
        isCard: React.PropTypes.bool.isRequired,
        isDeck: React.PropTypes.bool.isRequired,
        isReview: React.PropTypes.bool.isRequired,
        isStash: React.PropTypes.bool.isRequired
    },

    onClickDecks(event) {
        event.preventDefault();
        event.stopPropagation();

        this.props.store.invoke(toDeck);
    },

    onClickCards(event) {
        event.preventDefault();
        event.stopPropagation();

        this.props.store.invoke(toDeckCardsList);
    },

    onClickReview(event) {
        event.preventDefault();
        event.stopPropagation();

        this.props.store.invoke(toDeckReview);
    },

    onClickStashes(event) {
        event.preventDefault();
        event.stopPropagation();

        this.props.store.invoke(toStash);
    },

    render() {

        const {isCard, isDeck, isReview, isStash} = this.props;

        return (
            <div className="row">
                <div className="col-sm-6">
                    <div className="btn-group p-b pull-left" role="group" aria-label="Basic example">
                        <button
                            type="button"
                            className={classNames('btn', {'btn-primary': isDeck, 'btn-secondary': !isDeck})}
                            onClick={this.onClickDecks}>{"Decks"}</button>
                        <button
                            type="button"
                            className={classNames('btn', {'btn-primary': isCard, 'btn-secondary': !isCard})}
                            onClick={this.onClickCards}>{"Cards"}</button>
                        <button
                            type="button"
                            className={classNames('btn', {'btn-primary': isReview, 'btn-secondary': !isReview})}
                            onClick={this.onClickReview}>{"Review"}</button>
                    </div>
                </div>
                <div className="col-sm-6">
                    <div className="btn-group p-b pull-right" role="group" aria-label="Basic example">
                        <button
                            type="button"
                            className={classNames('btn', {'btn-primary': isStash, 'btn-secondary': !isStash})}
                            onClick={this.onClickStashes}>{"Stashes"}</button>
                    </div>
                </div>
            </div>
        );
    }
});

module.exports = orwell(SubNav, {

    watchCursors(props, manual, context) {

        const state = context.store.state();

        return [
            state.cursor(paths.dashboard.view)
        ];
    },

    assignNewProps(props, context) {

        const store = context.store;
        const state = store.state();

        const currentView = state.cursor(paths.dashboard.view).deref();

        const isCard = currentView === dashboard.view.cards;
        const isDeck = currentView === dashboard.view.decks;
        const isReview = currentView === dashboard.view.review;
        const isStash = currentView === dashboard.view.stash;

        return {
            store: store,
            isCard: isCard,
            isDeck: isDeck,
            isReview: isReview,
            isStash: isStash
        };
    }

}).inject({
    contextTypes: {
        store: React.PropTypes.object.isRequired
    }
});

