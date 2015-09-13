const React = require('react');
const orwell = require('orwell');
const Immutable = require('immutable');
const {navigateParentDeck} = require('store/decks');

const BreadcrumbChild = React.createClass({

    propTypes: {
        deck: React.PropTypes.instanceOf(Immutable.Map).isRequired,
        store: React.PropTypes.object.isRequired,
        active: React.PropTypes.bool
    },

    getDefaultProps() {
        return {
            active: false
        };
    },

    onClick(event) {
        event.preventDefault();
        event.stopPropagation();

        const {store, deck} = this.props;

        store.dispatch(navigateParentDeck, deck);
    },

    render() {

        const {deck} = this.props;

        const name = deck.get('name');
        const href = `#${deck.get('id')}/${name}`;

        if(this.props.active) {
            return (
                <li className="active">{name}</li>
            );
        }

        return (
            <li>
                <a key="a" href={href} onClick={this.onClick}>{name}</a>
            </li>
        );
    }
});

module.exports = orwell(BreadcrumbChild, {

    assignNewProps(props, context) {
        return {
            store: context.store
        };
    }
}).inject({
    contextTypes: {
        store: React.PropTypes.object.isRequired
    }
});
